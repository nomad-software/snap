// Package.
package action

// Imports.
import "fmt"
import "log"
import "snap/database"

// Initialise a datbase to be managed by snap. If the snap config database does 
// not exist when this function is called then it is automatically created.
func InitialiseDatabase(name string, recursion int) {

	if !database.ConfigDatabaseExists() {

		if recursion > 0 {
			// This check is here to make sure we don't enter a infinite loop 
			// trying to create the snap config database. This should never 
			// happen as creation errors are handled elsewhere but lets handle 
			// it just in case.
			log.Fatalln("Snap config database could not be created.")
		}
		log.Println("Snap config database does not exist")
		database.CreateConfigDatabase()

		// Use a recursion number greater than 0 to signal we are retrying this 
		// function with a newly created config database.
		InitialiseDatabase(name, 1)
		return
	}

	log.Println(fmt.Sprintf("Initialising '%s' database for snap managment", name))
}
