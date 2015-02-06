// Package.
package database

// Imports.
import "database/sql"
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
