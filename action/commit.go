// Package.
package action

// Imports.
import "github.com/nomad-software/snap/config"
import "github.com/nomad-software/snap/database"
import "github.com/nomad-software/snap/sanitise"
import "log"
import "strings"

// Commit a new file containing schema updates to a managed database.
func CommitFile(databaseName string, file string, comment string) {

	database.AssertConfigDatabaseExists()
	database.AssertDatabaseExists(databaseName)

	validateSqlFileFormat(file)

	database.ValidateSchemaUpdate(databaseName, file)
	database.CreateNewRevision(databaseName, file, comment)

	log.Println("File committed successfully.")
}

// Validate that the passed SQL file is off the required format.
func validateSqlFileFormat(file string) {
	upFound     := false
	downFound   := false
	contents    := sanitise.ReadFile(file)
	lines       := strings.Split(contents, "\n")
	for _, line := range lines {
		if line == config.UP_SQL_START && downFound == false {
			upFound = true
		} else if line == config.DOWN_SQL_START {
			downFound = true
		}
	}
	if !(upFound && downFound) {
		log.Fatalf("File '%s' is not in the correct format.", file)
	}
}
