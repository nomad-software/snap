// Package.
package sanitise

// Imports.
import "io/ioutil"
import "log"
import "strings"

// Convert the SQL string line endings to unix format.
func ConvertToUnixLineEndings(sql string) (string) {
	sql = strings.Replace(sql, "\r\n", "\n", -1);
	sql = strings.Replace(sql, "\r", "\n", -1);
	return sql
}

// Read a file and return the contents.
func ReadFile(name string) (string) {
	bytes, err := ioutil.ReadFile(name)
	if err != nil {
		log.Fatalln(err)
	}
	contents := string(bytes);
	return ConvertToUnixLineEndings(contents)
}
