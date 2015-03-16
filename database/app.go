// Package.
package database

// Imports.
import "crypto/rand"
import "fmt"
import "github.com/nomad-software/snap/config"
import "github.com/nomad-software/snap/sanitise"
import "log"
import "strings"

// Check if a database exists.
func DatabaseExists(name string) (bool) {
	err := useDatabase(name)
	return (err == nil)
}

// Assert the a database exists. If not throw a fatal error.
func AssertDatabaseExists(name string) {
	if !DatabaseExists(name) {
		log.Fatalf("Database '%s' does not exist.\n", name)
	}
}

// Check that a database is being managed.
func databaseIsManaged(database string) (bool) {
	AssertUseConfigDatabase()
	query := `SELECT id.name
		FROM initialisedDatabases AS id
		WHERE id.name = ?
		LIMIT 1;`
	row, err := QueryRow(query, database)
	exitOnError(err, "Error occurred checking database '%s' is being managed.", database)
	return len(row) != 0
}

// Assert that a database is being managed. If not throw a fatal error.
func assertDatabaseIsManaged(database string) {
	if !databaseIsManaged(database) {
		log.Fatalf("Database '%s' is not currently being managed.\n", database)
	}
}

// Add a database to be managed.
func InitialiseDatabase(database string) {

	fullSql := GenerateSchema(database)

	AssertUseConfigDatabase()
	StartTransaction()

		query := `INSERT INTO initialisedDatabases
			(name, currentSchemaRevision)
			VALUES (?, 1);`

		insertId, err := InsertRow(query, database)
		exitOnError(err, "Database '%s' is already being managed.", database)

		query = `INSERT INTO revisions
			(databaseId, revision, upSql, downSql, fullSql, comment, author)
			VALUES (?, 1, NULL, NULL, ?, "Database initialised.", ?);`

		_, err = InsertRow(query, insertId, fullSql, config.GetConfig().Identity)
		exitOnError(err, "Database '%s' is already being managed.", database)

	Commit()
}

// An initialised database type.
type database struct {
	Name string
	Revision string
	Date string
}

// A collection of initialised databases.
type databaseList []database

// Return a tabbed output string for writing using a tabbed writer.
func (this database) TabbedString() (string) {
	return fmt.Sprintf("%s\t%s\t%s", this.Name, this.Revision, this.Date)
}

// Get the length of the longest database name. This is to aid formatting a 
// command line ascii display.
func (this databaseList) LengthOfLongestName() (maxLength int) {
	// Set the default to the length of the ascii display's column heading i.e. 
	// 'Database'.
	maxLength = 8
	for _, entry := range this {
		if len(entry.Name) > maxLength {
			maxLength = len(entry.Name)
		}
	}
	return
}

// List all managed databases.
func GetManagedDatabaseList() (list databaseList) {

	AssertUseConfigDatabase()

	query := `SELECT id.name,
		MAX(r.revision) AS revision,
		id.dateInitialised
		FROM initialisedDatabases AS id
		INNER JOIN revisions AS r ON r.databaseId = id.id
		GROUP BY r.databaseId
		ORDER BY id.dateInitialised ASC;`

	rows, err := Query(query)
	exitOnError(err, "Can not retrieve list of managed databases.")

	list = make([]database, 0)
	for _, row := range rows {
		list = append(list, database{row.Str(0), row.Str(1), row.Str(2)})
	}
	return;
}

// A log entry type.
type logEntry struct {
	Revision string
	Comment string
	Author string
	Date string
}

// A collection of log entries.
type logEntries []logEntry

// Get log entries for the passed database.
func GetLogEntries(database string) (log logEntries) {

	assertDatabaseIsManaged(database)
	AssertUseConfigDatabase()

	query := `SELECT
		r.revision,
		r.comment,
		r.author,
		r.dateApplied
		FROM initialisedDatabases AS id
		INNER JOIN revisions AS r ON r.databaseId = id.id
		WHERE id.name = ?
		ORDER BY r.revision DESC;`

	rows, err := Query(query, database)
	exitOnError(err, "Can not retrieve log entries for database '%s'.", database)

	log = make([]logEntry, 0)
	for _, row := range rows {
		log = append(log, logEntry{row.Str(0), row.Str(1), row.Str(2), row.Str(3)})
	}
	return;
}

// Get the maximum revision of the passed database.
func GetHeadRevision(database string) (uint64) {

	assertDatabaseIsManaged(database)
	AssertUseConfigDatabase()

	query := `SELECT
		MAX(r.revision)
		FROM initialisedDatabases AS id
		INNER JOIN revisions AS r ON r.databaseId = id.id
		WHERE id.name = ?
		GROUP BY r.databaseId
		LIMIT 1;`

	row, err := QueryRow(query, database)
	exitOnError(err, "Can not retrieve latest revision for database '%s'.", database)

	return row.Uint64(0)
}

// Get the current schema revision of the passed database. The value returned 
// will be different to the latest revision if the database has been rolled 
// back using the update command.
func GetCurrentSchemaRevision(database string) (uint64) {

	assertDatabaseIsManaged(database)
	AssertUseConfigDatabase()

	query := `SELECT
		id.currentSchemaRevision
		FROM initialisedDatabases AS id
		WHERE id.name = ?
		LIMIT 1;`

	row, err := QueryRow(query, database)
	exitOnError(err, "Can not retrieve schema revision for database '%s'.", database)

	return row.Uint64(0)
}

// Set the current schema revision of the passed database.
func setCurrentSchemaRevision(database string, revision uint64) {
	assertDatabaseIsManaged(database)
	AssertUseConfigDatabase()

	query := `UPDATE initialisedDatabases AS id
		SET id.currentSchemaRevision = ?
		WHERE id.name = ?
		LIMIT 1;`

	err := Exec(query, revision, database)
	exitOnError(err, "Error occurred while setting the current schema revision for database '%s'.", database)
}

// Return the update SQL for the database and revision passed. This function 
// defaults to the full SQL if the update SQL doesn't exist. This is because 
// when initialising a database (revision 1) only the full SQL is available.
func GetUpdateSql(database string, revision uint64) (upSql string) {

	assertDatabaseIsManaged(database)
	AssertUseConfigDatabase()

	query := `SELECT
		COALESCE(r.upSql, r.fullSql)
		FROM initialisedDatabases AS id
		INNER JOIN revisions AS r ON r.databaseId = id.id
		WHERE id.name = ?
		AND r.revision = ?
		LIMIT 1;`

	row, err := QueryRow(query, database, revision)
	exitOnError(err, "Can not retrieve update SQL for database '%s' at revision '%d'.", database, revision)

	if len(row) > 0 {
		upSql = row.Str(0)
	}
	return
}

// Return the down SQL for the database and revision passed.
func GetDownSql(database string, revision uint64) (downSql string) {

	assertDatabaseIsManaged(database)
	AssertUseConfigDatabase()

	query := `SELECT
		r.downSql
		FROM initialisedDatabases AS id
		INNER JOIN revisions AS r ON r.databaseId = id.id
		WHERE id.name = ?
		AND r.revision = ?
		LIMIT 1;`

	row, err := QueryRow(query, database, revision)
	exitOnError(err, "Can not retrieve down SQL for database '%s' at revision '%d'.", database, revision)

	if len(row) > 0 {
		downSql = row.Str(0)
	}
	return
}

// Return the full SQL for the database and revision passed.
func GetSchema(database string, revision uint64) (sql string) {

	assertDatabaseIsManaged(database)
	AssertUseConfigDatabase()

	query := `SELECT
		r.fullSql
		FROM initialisedDatabases AS id
		INNER JOIN revisions AS r ON r.databaseId = id.id
		WHERE id.name = ?
		AND r.revision = ?
		LIMIT 1;`

	row, err := QueryRow(query, database, revision)
	exitOnError(err, "Can not retrieve full SQL for database '%s' at revision '%d'.", database, revision)

	if len(row) > 0 {
		sql = row.Str(0)
	}
	return
}

// Copy a full source database (sans data) to a new destination at a particular 
// source revision.
func CopyDatabase(source string, destination string, revision uint64) {

	assertDatabaseIsManaged(source)
	charSet, collation := GetDatabaseEncoding(source)
	SetConnectionEncoding(charSet, collation)

	err := createDatabase(destination, charSet, collation)
	exitOnError(err, "Can not create new database '%s'.", destination)

	sql := GetSchema(source, revision)
	sql  = sanitise.SanitiseSql(sql)

	assertUseDatabase(destination)
	err = ExecMulti(sql)
	exitOnError(err, "Can not copy schema to new database '%s'.", destination)
}

// Validate that the schema file updates then correctly reverses any changes made.
func ValidateSchemaUpdate(database string, file string) {
	assertDatabaseIsManaged(database)

	temp     := generateTempDatabaseName()
	revision := GetHeadRevision(database)
	CopyDatabase(database, temp, revision)

	sql := sanitise.ReadFile(file)
	sql  = sanitise.SanitiseSql(sql)

	assertUseDatabase(temp)
	err := ExecMulti(sql)
	exitOnError(err, "Error occurred applying file to current schema.")

	currentStructure := GetSchema(database, revision)
	updatedStructure := GenerateSchema(temp)
	updatedStructure  = strings.Replace(updatedStructure, temp, database, -1)

	deleteTempDatabases()

	if currentStructure != updatedStructure {
		log.Fatalln("File not commited because it doesn't correctly reverse any contained updates.")
	}
}

// Create a new revision for a managed database. This function applies the file 
// and creates the new revision in the database.
func CreateNewRevision(database string, file string, comment string) {
	assertDatabaseIsManaged(database)

	sql := sanitise.ReadFile(file)
	sql  = sanitise.SanitiseSql(sql)

	databaseId     := getDatabaseId(database)
	revision       := GetHeadRevision(database) + 1
	upSql, downSql := splitSqlFile(sql)
	author         := config.GetConfig().Identity

	StartTransaction()

		applyUpdateToDatabase(database, upSql)
		fullSql := GenerateSchema(database)

		AssertUseConfigDatabase()
		query := `INSERT INTO revisions
			(databaseId, revision, upSql, downSql, fullSql, comment, author)
			VALUES (?, ?, ?, ?, ?, ?, ?);`

		_, err := InsertRow(query, databaseId, revision, upSql, downSql, fullSql, comment, author)
		exitOnError(err, "Error occurred while creating a new revision for database '%s'.", database)

		setCurrentSchemaRevision(database, revision)

	Commit()
}

// Get the id of a managed database.
func getDatabaseId(database string) (uint64) {
	assertDatabaseIsManaged(database)
	AssertUseConfigDatabase()
	query := `SELECT id.id
		FROM initialisedDatabases AS id
		WHERE id.name = ?
		LIMIT 1;`
	row, err := QueryRow(query, database)
	exitOnError(err, "Error occurred while retrieving database '%s' id.", database)
	return row.Uint64(0)
}

// Apply the update SQL to a database.
func applyUpdateToDatabase(database string, sql string) {
	assertDatabaseIsManaged(database)
	assertUseDatabase(database)
	err := ExecMulti(sql)
	exitOnError(err, "Error occurred while modifying database '%s' schema.", database)
}

// Split the update SQL into the up and down sections.
// This function assumes the SQL has been validated before hand.
func splitSqlFile(sql string) (upSql string, downSql string) {
	upLines   := make([]string, 0)
	downLines := make([]string, 0)
	var output *[]string
	lines := strings.Split(sql, "\n")
	for _, line := range lines {
		if line == config.UP_SQL_START {
			output = &upLines
			continue
		} else if line == config.DOWN_SQL_START {
			output = &downLines
			continue
		}
		*output = append(*output, line)
	}
	upSql   = strings.Join(upLines, "\n");
	downSql = strings.Join(downLines, "\n");
	return
}

// Generate a random name for a temporary database.
func generateTempDatabaseName() (string) {
	bytes := make([]byte, 4)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatalln("Error occurred generating a temporary database name.")
	}
	name := fmt.Sprintf("snap_%X", bytes)
	// Record the name to drop later to clean up.
	tempDatabases = append(tempDatabases, name)
	return name;
}

// Update the schema of a managed database to a previously committed revision.
func UpdateSchemaToRevision(database string, target uint64) {
	assertDatabaseIsManaged(database)
	revision := GetCurrentSchemaRevision(database)
	if target > revision {
		for revision < target {
			revision++
			forwardSchema(database, revision)
		}
	} else {
		for revision > target {
			reverseSchema(database, revision)
			revision--
		}
	}
}

// Foward the schema to a stored revision.
func forwardSchema(database string, target uint64) {
	assertDatabaseIsManaged(database)
	sql := GetUpdateSql(database, target)
	assertUseDatabase(database)
	err := ExecMulti(sql)
	exitOnError(err, "Error occurred applying update SQL to database '%s'.", database)
	setCurrentSchemaRevision(database, target)
}

// Reverse the schema to a stored revision.
func reverseSchema(database string, target uint64) {
	assertDatabaseIsManaged(database)
	sql := GetDownSql(database, target)
	assertUseDatabase(database)
	err := ExecMulti(sql)
	exitOnError(err, "Error occurred applying down SQL to database '%s'.", database)
	target--
	setCurrentSchemaRevision(database, target)
}
