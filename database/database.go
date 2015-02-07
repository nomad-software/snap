// Package.
package database

// Imports.
import "database/sql"
import "fmt"
import "log"
import "snap/config"
import _ "github.com/go-sql-driver/mysql"

var db *sql.DB;
var rows *sql.Rows
var err error

// Establishes a connection to the database.
func OpenConnection(config config.Database) {
	db, err = sql.Open("mysql", config.String())
	if err != nil {
		log.Fatalln(err)
	}
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
	_, err = db.Query(fmt.Sprintf("USE %s;", name))
	return err;
}

// Drop a database.
func DropDatabase(name string) (error) {
	_, err = db.Query(fmt.Sprintf("DROP SCHEMA IF EXISTS %s;", name))
	return err;
}

// Create a database.
func CreateDatabase(name string) (error) {
	_, err = db.Query(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;", name))
	return err;
}
