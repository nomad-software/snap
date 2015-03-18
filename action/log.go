package action

// Imports.
import "fmt"
import "github.com/nomad-software/snap/database"
import "log"

// Show the commit log for the passed database.
func ShowLog(databaseName string) {

	database.AssertConfigDatabaseExists()
	database.AssertDatabaseExists(databaseName)

	logEntries := database.GetLogEntries(databaseName)

	if len(logEntries) > 0 {
		for _, entry := range logEntries {
			fmt.Printf("Revision: %s\n", entry.Revision)
			fmt.Printf("Author: %s\n", entry.Author)
			fmt.Printf("Date: %s\n", entry.Date)
			fmt.Println("")
			fmt.Printf("    %s\n", entry.Comment)
			fmt.Println("")
		}
	} else {
		log.Printf("No log entries found for database '%s'.\n", databaseName)
	}
}
