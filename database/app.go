// Package.
package database

// Imports.
import "fmt"
import "github.com/nomad-software/snap/config"
import "log"

// Check that a database is being managed.
func databaseIsManaged(databaseName string) (bool) {
	UseConfigDatabase()
	query := `SELECT id.name
		FROM initialisedDatabases AS id
		WHERE id.name = ?
		LIMIT 1;`
	row, err := QueryRow(query, databaseName)
	exitOnError(err, fmt.Sprintf("Error occurred checking database '%s' is being managed.", databaseName))
	return len(row) != 0
}

// Assert that a database is being managed. If not throw a fatal error.
func assertDatabaseIsManaged(databaseName string) {
	if !databaseIsManaged(databaseName) {
		log.Fatalf("Database '%s' is not currently being managed.\n", databaseName)
	}
}

// Add a database to be managed.
func InitialiseDatabase(databaseName string) {

	fullSql := GenerateFullSql(databaseName)

	UseConfigDatabase()
	StartTransaction()

		insertId, err := InsertRow("INSERT INTO initialisedDatabases (name) VALUES (?)", databaseName)
		exitOnError(err, fmt.Sprintf("Database '%s' is already being managed.", databaseName))

		query := `INSERT INTO revisions
			(databaseId, revision, upSql, downSql, fullSql, comment, author)
			VALUES (?, 1, NULL, NULL, ?, "Database initialised.", ?);`

		_, err = InsertRow(query, insertId, fullSql, config.GetConfig().Identity)
		exitOnError(err, fmt.Sprintf("Database '%s' is already being managed.", databaseName))

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

// Get the length of the longest database name.
func (this databaseList) LengthOfLongestName() (maxLength int) {
	// Set the default to the same as the length of this column's heading i.e. 'Database'.
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

	UseConfigDatabase()

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
func GetLogEntries(databaseName string) (log logEntries) {

	assertDatabaseIsManaged(databaseName)
	UseConfigDatabase()

	query := `SELECT
		r.revision,
		r.comment,
		r.author,
		r.dateApplied
		FROM initialisedDatabases AS id
		INNER JOIN revisions AS r ON r.databaseId = id.id
		WHERE id.name = ?
		ORDER BY r.revision DESC;`

	rows, err := Query(query, databaseName)
	exitOnError(err, fmt.Sprintf("Can not retrieve log entries for database '%s'.\n", databaseName))

	log = make([]logEntry, 0)
	for _, row := range rows {
		log = append(log, logEntry{row.Str(0), row.Str(1), row.Str(2), row.Str(3)})
	}
	return;
}

// Get the current revision of the passed database.
func GetCurrentRevision(databaseName string) (uint64) {

	assertDatabaseIsManaged(databaseName)
	UseConfigDatabase()

	query := `SELECT
		MAX(r.revision)
		FROM initialisedDatabases AS id
		INNER JOIN revisions AS r ON r.databaseId = id.id
		WHERE id.name = ?
		GROUP BY r.databaseId
		LIMIT 1;`

	row, err := QueryRow(query, databaseName)
	exitOnError(err, fmt.Sprintf("Can not retrieve current revision for database '%s'.\n", databaseName))

	return row.Uint64(0)
}

// Return the update SQL for the database and revision passed. This function 
// defaults to the full SQL if the update SQL doesn't exist. This is because 
// when initialising a database (revision 1) only the full SQL is available.
func GetUpdateSql(databaseName string, revision uint64) (upSql string) {

	assertDatabaseIsManaged(databaseName)
	UseConfigDatabase()

	query := `SELECT
		COALESCE(r.upSql, r.fullSql)
		FROM initialisedDatabases AS id
		INNER JOIN revisions AS r ON r.databaseId = id.id
		WHERE id.name = ?
		AND r.revision = ?
		LIMIT 1;`

	row, err := QueryRow(query, databaseName, revision)
	exitOnError(err, fmt.Sprintf("Can not retrieve update SQL for database '%s' at revision '%d'.\n", databaseName, revision))

	if len(row) > 0 {
		upSql = row.Str(0)
	}
	return
}

// Return the full SQL for the database and revision passed.
func GetFullSql(databaseName string, revision uint64) (fullSql string) {

	assertDatabaseIsManaged(databaseName)
	UseConfigDatabase()

	query := `SELECT
		r.fullSql
		FROM initialisedDatabases AS id
		INNER JOIN revisions AS r ON r.databaseId = id.id
		WHERE id.name = ?
		AND r.revision = ?
		LIMIT 1;`

	row, err := QueryRow(query, databaseName, revision)
	exitOnError(err, fmt.Sprintf("Can not retrieve full SQL for database '%s' at revision '%d'.\n", databaseName, revision))

	if len(row) > 0 {
		fullSql = row.Str(0)
	}
	return
}

// Copy a full database to a destination at a particular revision.
func CopyDatabase(sourceDatabaseName string, destinationDatabaseName string, revision uint64) {

	charSet, collation := GetDatabaseEncoding(sourceDatabaseName)
	SetConnectionEncoding(charSet, collation)

	log.Printf("Creating database '%s' with charset '%s' and collation '%s'.", destinationDatabaseName, charSet, collation)

	err := createDatabase(destinationDatabaseName, charSet, collation)
	exitOnError(err, fmt.Sprintf("Can not create new database '%s'.\n", destinationDatabaseName))

	fullSql := GetFullSql(sourceDatabaseName, revision)
	fullSql  = sanitise(fullSql)

	AssertUseDatabase(destinationDatabaseName)
	log.Printf("Copying schema from '%s' to '%s'.", sourceDatabaseName, destinationDatabaseName)

	err = ExecMulti(fullSql)
	exitOnError(err, fmt.Sprintf("Can not copy schema to new database '%s'.\n", destinationDatabaseName))

	log.Println("Success.")
}
