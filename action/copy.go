// Package.
package action

// Imports.
import "github.com/nomad-software/snap/database"
import "log"

// Copy a full database to a destination at a particular revision.
func CopyDatabase(source string, destination string, revision uint64) {

	database.AssertConfigDatabaseExists()
	database.AssertDatabaseExists(source)

	if database.DatabaseExists(destination) {
		log.Fatalf("Database '%s' already exists.\n", destination)
	}

	currentRevision := database.GetCurrentRevision(source)

	if revision > currentRevision {
		log.Fatalf("Database '%s' does not have a revision '%d'.\n", source, revision)
	}

	if revision == 0 {
		revision = currentRevision
	}

	database.CopyDatabase(source, destination, revision)

	log.Println("Database copied successfully.")
}
