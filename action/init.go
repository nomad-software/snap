// Package.
package action

// Imports.
import "fmt"
import "log"
import "github.com/nomad-software/snap/database"

// Initialise a datbase to be managed by snap.
func InitialiseDatabase(name string) {

	database.AssertConfigDatabaseExists()
	database.AssertDatabaseExists(name)

	log.Println(fmt.Sprintf("Initialising '%s' database for snap managment", name))
	database.InitialiseDatabase(name)
}
