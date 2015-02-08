// Package.
package action

// Imports.
import "fmt"
import "log"
import "snap/database"

// Check for the existance of the config database. If it doesn't exist, create 
// it.
func checkConfigDbExists() {
	if !database.ConfigDatabaseExists() {
		log.Println("Snap config database does not exist")
		database.CreateConfigDatabase()
	}
}

// Initialise a datbase to be managed by snap.
func InitialiseDatabase(name string) {

	checkConfigDbExists()

	if database.DatabaseExists(name) {
		log.Println(fmt.Sprintf("Initialising '%s' database for snap managment", name))
		database.InitialiseDatabase(name)
	} else {
		log.Println(fmt.Sprintf("Database '%s' does not exist.", name))
	}
}
