// Package.
package database

// Imports.
import "database/sql"
import "fmt"

// Add a database to be managed.
func InitialiseDatabase(name string) {

	var transaction *sql.Tx
	var result sql.Result
	var err error
	var insertId int64

	UseConfigDatabase()

	transaction, err = db.Begin()
	ExitOnError(err, fmt.Sprintf("Error occured while trying to manage database '%s'.", name))

	result, err = transaction.Exec("INSERT INTO initialisedDatabases (name) VALUES (?);", name)
	ExitOnError(err, fmt.Sprintf("Database '%s' is already being managed.", name))

	if result != nil {
		insertId, err = result.LastInsertId()

		query := `INSERT INTO revisions
			(databaseId, revision, upSql, downSql, fullSql, comment)
			VALUES (?, 1, "", "", "", "Database initialised.");`

		_, err = transaction.Exec(query, insertId)
		ExitOnError(err, fmt.Sprintf("Database '%s' is already being managed.", name))
	}

	err = transaction.Commit()
	ExitOnError(err, fmt.Sprintf("Error occured while trying to manage database '%s'.", name))
}
