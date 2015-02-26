// Package.
package action

// Imports.
import "fmt"
import "github.com/nomad-software/snap/database"
import "log"

// Show a managed database's full SQL at a particular revision.
func ShowFullSql(databaseName string, revision uint64) {

	database.AssertConfigDatabaseExists()
	database.AssertDatabaseExists(databaseName)

	currentRevision := database.GetCurrentRevision(databaseName)

	if revision > currentRevision {
		log.Fatalf("Database '%s' does not have a revision '%d'.\n", databaseName, revision)
	}

	if revision == 0 {
		revision = currentRevision
	}

	sql := database.GetFullSql(databaseName, revision)
	fmt.Println(sql)
}
