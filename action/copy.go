// Package.
package action

// Imports.
import "github.com/nomad-software/snap/database"
import "log"

// Copy a full database to a destination at a particular revision.
func CopyDatabase(sourceDatabaseName string, destinationDatabaseName string, revision uint64) {

	database.AssertConfigDatabaseExists()
	database.AssertDatabaseExists(sourceDatabaseName)

	if database.DatabaseExists(destinationDatabaseName) {
		log.Fatalf("Database '%s' already exists.\n", destinationDatabaseName)
	}

	currentRevision := database.GetCurrentRevision(sourceDatabaseName)

	if revision > currentRevision {
		log.Fatalf("Database '%s' does not have a revision '%d'.\n", sourceDatabaseName, revision)
	}

	if revision == 0 {
		revision = currentRevision
	}

	database.CopyDatabase(sourceDatabaseName, destinationDatabaseName, revision)
}
