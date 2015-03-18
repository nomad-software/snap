package action

// Imports.
import "fmt"
import "github.com/nomad-software/snap/database"
import "log"

// Show a managed database's full SQL at a particular revision.
func ShowFullSql(databaseName string, revision uint64) {

	database.AssertConfigDatabaseExists()
	database.AssertDatabaseExists(databaseName)

	head := database.GetHeadRevision(databaseName)

	if revision > head {
		log.Fatalf("Database '%s' does not have a revision '%d'.\n", databaseName, revision)
	}

	if revision <= 0 {
		revision = head
	}

	fullSql := database.GetSchema(databaseName, revision)
	fmt.Println(fullSql)
}
