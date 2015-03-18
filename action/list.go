package action

// Imports.
import "fmt"
import "github.com/nomad-software/snap/database"
import "log"
import "os"
import "strings"
import "text/tabwriter"

// List all managed databases.
func ListManagedDatabases() {

	database.AssertConfigDatabaseExists()

	list := database.GetManagedDatabaseList()

	if len(list) > 0 {

		writer := tabwriter.NewWriter(os.Stdout, 8, 4, 1, ' ', 0)
		fmt.Fprintln(writer, "Database\tRevision\tInitialised")

		firstColumnLine := strings.Repeat("-", list.LengthOfLongestName())

		fmt.Fprintf(writer, "%s\t--------\t-------------------\n", firstColumnLine)

		for _, entry := range list {
			fmt.Fprintln(writer, entry.TabbedString())
		}
		writer.Flush()
	} else {
		log.Println("No databases are currently being managed.")
	}
}
