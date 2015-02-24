// Package.
package database

// Imports.
import "fmt"

// Add a database to be managed.
func InitialiseDatabase(name string) {

	fullSql := DumpDatabase(name)

	UseConfigDatabase()
	StartTransaction()

		insertId, err := InsertRow("INSERT INTO initialisedDatabases (name) VALUES (?)", name)
		ExitOnError(err, fmt.Sprintf("Database '%s' is already being managed.", name))

		query := `INSERT INTO revisions
			(databaseId, revision, upSql, downSql, fullSql, comment)
			VALUES (?, 1, NULL, NULL, ?, "Database initialised.");`

		_, err = InsertRow(query, insertId, fullSql)
		ExitOnError(err, fmt.Sprintf("Database '%s' is already being managed.", name))

	Commit()
}

// Type returned by ListManagedDatabases function.
type database struct {
	Name string
	Revision string
	Date string
}

// A collection of databases.
type databaseList []database

// Return a tabbed output string for writing using a tabbed writer.
func (this database) TabbedString() (string) {
	return fmt.Sprintf("%s\t%s\t%s", this.Name, this.Revision, this.Date)
}

// Get the length of the longest database name.
func (this databaseList) LengthOfLongestName() (maxLength int) {
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
	ExitOnError(err, "Can not retrieve list of managed databases.")

	list = make([]database, 0)
	for _, row := range rows {
		list = append(list, database{row.Str(0), row.Str(1), row.Str(2)})
	}
	return;
}
