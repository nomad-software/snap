package action

// Imports.
import "github.com/nomad-software/snap/database"
import "log"

// Show a managed database's full SQL at a particular revision.
func UpdateSchemaToRevision(databaseName string, target uint64) {

	database.AssertConfigDatabaseExists()
	database.AssertDatabaseExists(databaseName)

	head    := database.GetHeadRevision(databaseName)
	current := database.GetCurrentSchemaRevision(databaseName)

	if target <= 0 {
		target = head
	}

	if target > head {
		log.Fatalf("Database '%s' does not have a revision '%d'.\n", databaseName, target)
	} else if target == current {
		log.Fatalf("Database '%s' is already at target revision '%d'.\n", databaseName, target)
	}

	database.UpdateSchemaToRevision(databaseName, target)
}
