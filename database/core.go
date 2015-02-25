// Package.
package database

// Imports.
import "fmt"
import "github.com/nomad-software/snap/config"
import "github.com/ziutek/mymysql/mysql"
import "log"
import _ "github.com/ziutek/mymysql/native"

// Package database struct.
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

// Handle a fatal error that will halt program execution. Rollback any 
// transaction that is pending.
func ExitOnError(err error, messages ...interface{}) {
	if err != nil {
		Rollback()
		log.Println(err)
		log.Fatalln(messages...)
	}
}

// Establishes a connection to the database.
func Open(config config.Config) {
	protocol      := config.Database.Protocol
	localAddress  := ""
	remoteAddress := config.Database.Host + ":" + config.Database.Port
	user          := config.Database.User
	password      := config.Database.Password
	database      := ""
	_db := mysql.New(protocol, localAddress, remoteAddress, user, password, database)
	err := _db.Connect()
	ExitOnError(err, "Database connection could not be established.")
	db = _db
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

// Execute an unsafe statement not expecting results. This function is used 
// very similarly to the fmt.Sprintf function and all format specifiers are 
// supported. Escaping of parameters is handled in the wrapped library.
func ExecUnsafe(sql string, params ...interface{}) (err error) {
	_, _, err = db.Query(sql, params...)
	return
}

// Execute a multi-statement query expecting no results. This is especially 
// useful for executing many SQL statements in one go, such as applying DDL's 
// to an existing schema.
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

// Execute a prepared statement expecting multiple results.
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

// Execute an unsafe statement expecting multiple results. This function is 
// used very similarly to the fmt.Sprintf function and all format specifiers are 
// supported. Escaping of parameters is handled in the wrapped library.
func QueryUnsafe(sql string, params ...interface{}) (rows []mysql.Row, err error) {
	rows, _, err = db.Query(sql, params...)
	return
}

// Execute a prepared statement expecting a single row result.
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

// Execute an unsafe statement expecting a single row result. This function is 
// used very similarly to the fmt.Sprintf function and all format specifiers are 
// supported. Escaping of parameters is handled in the wrapped library.
func QueryRowUnsafe(sql string, params ...interface{}) (row mysql.Row, err error) {
	rows, _, err := db.Query(sql, params...)
	if len(rows) > 0 {
		row = rows[0]
	}
	return
}

// Execute a prepared statement to insert a single row. The insert id is 
// returned.
func InsertRow(sql string, params ...interface{}) (insertId uint64, err error) {
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
	ExitOnError(err, "Error occurred starting transaction.")
	tx = _tx
}

// Commit a transaction.
func Commit() {
	if tx != nil && tx.IsValid() {
		err := tx.Commit()
		ExitOnError(err, "Error occurred committing transaction.")
	}
}

// Rollback a transaction.
func Rollback() {
	if tx != nil && tx.IsValid() {
		err := tx.Rollback()
		ExitOnError(err, "Error occurred rolling back transaction.")
	}
}

// Create a database.
func CreateDatabase(name string) (error) {
	err := ExecUnsafe("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;", name)
	return err
}

// Change the database to the one named in the name parameter.
func useDatabase(name string) (error) {
	err := ExecUnsafe("USE `%s`;", name)
	return err
}

// Assert the database can be used. If not throw a fatal error.
func AssertUseDatabase(name string) {
	err := useDatabase(name)
	ExitOnError(err, fmt.Sprintf("Can not use '%s' database.", name))
}

// Check if a database exists.
func DatabaseExists(name string) (bool) {
	err := useDatabase(name)
	return WasSuccessful(err)
}

// Assert the a database exists. If not throw a fatal error.
func AssertDatabaseExists(name string) {
	if !DatabaseExists(name) {
		log.Fatalf("Database '%s' does not exist.\n", name)
	}
}

// Drop a database.
func DropDatabase(name string) (error) {
	err := ExecUnsafe("DROP SCHEMA IF EXISTS `%s`;", name)
	return err
}
