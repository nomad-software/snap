// Package.
package action

// Imports.
import "fmt"
import "log"
import "github.com/nomad-software/snap/database"

// Initialise a datbase to be managed.
func InitialiseDatabase(databaseName string) {

	database.AssertConfigDatabaseExists()
	database.AssertDatabaseExists(databaseName)

	log.Println(fmt.Sprintf("Initialising '%s' database for managment", databaseName))
	database.InitialiseDatabase(databaseName)
}
