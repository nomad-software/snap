package action

// Imports.
import "fmt"
import "github.com/nomad-software/snap/database"
import "io/ioutil"
import "log"
import "os/exec"
import "strconv"
import "strings"

// Copy a full database to a destination at a particular revision.
func Diff(databaseName string, revisionString string) {

	database.AssertConfigDatabaseExists()
	database.AssertDatabaseExists(databaseName)

	from, to := parseRevisions(revisionString)
	if to == 0 {
		to = database.GetHeadRevision(databaseName)
	}
	if from > to {
		log.Fatalln("'From' revision cannot be greater than to revision.")
	}

	fromFile := fmt.Sprintf("/tmp/revision-%d", from)
	toFile   := fmt.Sprintf("/tmp/revision-%d", to)

	fromSql  := database.GetSchema(databaseName, from)
	toSql    := database.GetSchema(databaseName, to)

	writeFile(fromFile, fromSql)
	writeFile(toFile, toSql)

	output, err := exec.Command("diff", "-u", fromFile, toFile).CombinedOutput()
	if err != nil {
		switch err.(type) {
			case *exec.ExitError:
				// Diff returns an exit code of 1 if the files are different. 
				// So lets just skip that case.
			default:
				log.Fatalln(err)
		}
	}

	fmt.Println(string(output))
}

// Parse the revisions from the revision string.
func parseRevisions(revisionString string) (from uint64, to uint64) {
	var err error
	if strings.Contains(revisionString, "..") {
		revisions := strings.Split(revisionString, "..")
		if len(revisions) == 2 {
			from, err = strconv.ParseUint(revisions[0], 10, 64)
			if err != nil {
				log.Fatalf("'From' revision can not be recognised in '%s'.\n", revisionString)
			}
			to, err = strconv.ParseUint(revisions[1], 10, 64)
			if err != nil {
				log.Fatalf("'To' revision can not be recognised in '%s'.\n", revisionString)
			}
		} else {
			log.Fatalf("Revisions '%s' are not specified correctly.\n", revisionString)
		}
	} else {
		from, err = strconv.ParseUint(revisionString, 10, 64)
		if err != nil {
			log.Fatalf("Revision '%s' is not specified correctly.\n", revisionString)
		}
	}
	return
}

// Write text to a file.
func writeFile(file string, text string) {
	err := ioutil.WriteFile(file, []byte(text), 0644)
	if err != nil {
		log.Fatalln("Error writing to temporary file '%s'", file)
	}
}
