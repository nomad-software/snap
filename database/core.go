// Package.
package database

// Imports.
import "fmt"
import "github.com/ziutek/mymysql/mysql"
import "log"
import "snap/config"
import _ "github.com/ziutek/mymysql/native"

// Global database struct.
var db mysql.Conn
var tx mysql.Transaction

// Check if an error occurred. If it did print the error message and return 
// false. Return true if there was no error.
func WasSuccessful(err error) (bool) {
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// Handle a fatal error that will halt program execution.
func ExitOnError(err error, messages ...interface{}) {
	if err != nil {
		log.Println(err)
		log.Fatalln(messages...)
	}
}

// Establishes a connection to the database.
func Open(config config.Database) {
	_db := mysql.New(config.Protocol, "", config.Host + ":" + config.Port, config.User, config.Password, "")
	err := _db.Connect()
	db = _db
	ExitOnError(err, "Database connection could not be established.")
}

// Close the datbase connection.
func Close() {
	db.Close()
}

// Execute a prepared statement not expecting results.
func Exec(sql string, params ...interface{}) (err error) {
	statement, err := db.Prepare(sql)
	if err != nil {
		return
	}
	_, err = statement.Run(params...)
	return
}

// Execute an unsafe statement not expecting results, where all string 
// parameters should be escaped. The format of the SQL string should be the 
// same as the fmt.Sprintf function.
func ExecUnsafe(sql string, params ...interface{}) (err error) {
	escaped := make([]interface{}, 0)
	for _, param := range params {
		value, ok := param.(string)
		if ok {
			escaped = append(escaped, db.Escape(value))
		} else {
			escaped = append(escaped, param)
		}
	}
	_, _, err = db.Query(sql, escaped...)
	return err
}

// Execute a multi-statement query expecting no results. This is especially 
// useful for executing many SQL statements in one go.
func ExecMulti(sql string) (error) {
	result, err := db.Start(sql)
	if err != nil {
		return err
	}
	for result.MoreResults() {
		result, err = result.NextResult()
		if err != nil {
			return err
		}
		result.End()
	}
	return err
}

// Execute a prepared statement expecting results.
func Query(sql string, params ...interface{}) (rows []mysql.Row, err error) {
	statement, err := db.Prepare(sql)
	if err != nil {
		return
	}
	result, err := statement.Run(params...)
	if err != nil {
		return
	}
	rows, err = result.GetRows()
	return
}

// Execute a prepared statement expecting results.
func QueryRow(sql string, params ...interface{}) (row mysql.Row, err error) {
	statement, err := db.Prepare(sql)
	if err != nil {
		return
	}
	result, err := statement.Run(params...)
	if err != nil {
		return
	}
	row, err = result.GetFirstRow()
	return
}

// Execute a prepared statement to insert data. The insert id is returned.
func Insert(sql string, params ...interface{}) (insertId uint64, err error) {
	statement, err := db.Prepare(sql)
	if err != nil {
		return
	}
	result, err := statement.Run(params...)
	if err != nil {
		return
	}
	insertId = result.InsertId()
	return
}

// Start a transaction.
func StartTransaction() {
	_tx, err := db.Begin()
	tx = _tx
	ExitOnError(err, "Error occurred starting transaction.")
}

// Commit a transaction.
func Commit() {
	err := tx.Commit()
	ExitOnError(err, "Error occurred committing transaction.")
}

// Rollback a transaction.
func Rollback() {
	err := tx.Rollback()
	ExitOnError(err, "Error occurred rolling back transaction.")
}

// Create a database.
func CreateDatabase(name string) (error) {
	err := ExecUnsafe("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;", name)
	return err
}

// Drop a database.
func DropDatabase(name string) (error) {
	err := ExecUnsafe("DROP SCHEMA IF EXISTS %s;", name)
	return err
}
// Check if a database exists.
func DatabaseExists(name string) (bool) {
	err := UseDatabase(name)
	return WasSuccessful(err)
}

// Assert the a database exists. If not throw a fatal error.
func AssertDatabaseExists(name string) {
	if !DatabaseExists(name) {
		log.Println(fmt.Sprintf("Database '%s' does not exist.", name))
	}
}

// Change the database to the one named in the name parameter.
func UseDatabase(name string) (error) {
	err := ExecUnsafe("USE %s;", name)
	return err
}

// Asser the database can be used. If not throw a fatal error.
func AssertUseDatabase(name string) {
	err := UseDatabase(name)
	ExitOnError(err, fmt.Sprintf("Can not use '%s' database.", name))
}
