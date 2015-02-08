// Package.
package database

// Imports.
import "database/sql"
import "fmt"
import "log"
import "snap/config"
import _ "github.com/go-sql-driver/mysql"

// Global database struct.
var db *sql.DB

// Check if an error occurred. If it did print the error message. Return true 
// or false depending on whether an error occurred.
func WasSuccessful(err error) (bool) {
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// Handle a fatal error that will halt program execution.
func ExitOnError(err error, message string) {
	if err != nil {
		log.Fatalln(message, fmt.Sprintf("(%s)", err))
	}
}

// Establishes a connection to the database.
func OpenConnection(config config.Database) {
	_db, err := sql.Open("mysql", config.String())
	db = _db
	ExitOnError(err, "Database connection could not be established.")
}

// Close the datbase connection.
func Close() {
	db.Close()
}

// Check if a database exists.
func DatabaseExists(name string) (bool) {
	err := UseDatabase(name)
	return err == nil
}

// Change the database to the one named in the name parameter.
func UseDatabase(name string) (error) {
	_, err := db.Exec(fmt.Sprintf("USE %s;", name))
	return err
}

// Drop a database.
func DropDatabase(name string) (error) {
	_, err := db.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s;", name))
	return err
}

// Create a database.
func CreateDatabase(name string) (error) {
	_, err := db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;", name))
	return err
}
