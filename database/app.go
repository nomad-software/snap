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
