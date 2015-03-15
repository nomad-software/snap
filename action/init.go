// Package.
package action

// Imports.
import "log"
import "github.com/nomad-software/snap/database"

// Initialise a datbase to be managed.
func InitialiseDatabase(databaseName string) {

	database.AssertConfigDatabaseExists()
	database.AssertDatabaseExists(databaseName)

	log.Printf("Initialising '%s' database for managment\n", databaseName)
	database.InitialiseDatabase(databaseName)
	log.Println("Database initialised successfully.")
}
