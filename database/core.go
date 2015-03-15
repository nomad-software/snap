// Package.
package database

// Imports.
import "github.com/nomad-software/snap/config"
import "github.com/ziutek/mymysql/mysql"
import "log"
import _ "github.com/ziutek/mymysql/native"

// Package database struct.
var db mysql.Conn
var tx mysql.Transaction
var tempDatabases []string = make([]string, 0)

// Handle a fatal error that will halt program execution. Rollback any 
// transaction that is pending.
func exitOnError(err error, format string, values ...interface{}) {
	if err != nil {
		Rollback()
		deleteTempDatabases()
		log.Println(err)
		log.Fatalf(format + "\n", values...)
	}
}

// Delete any temporary database that have been created.
func deleteTempDatabases() {
	if len(tempDatabases) > 0 {
		for _, database := range tempDatabases {
			_ = dropDatabase(database)
		}
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
	exitOnError(err, "Database connection could not be established.")
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

// Set the connection encoding.
func SetConnectionEncoding(charSet string, collation string) {
	charSetQueries := []string{
		"SET character_set_client = ?",
		"SET character_set_results = ?",
		"SET character_set_connection = ?",
	}
	collationQueries := []string{
		"SET collation_connection = ?",
	}
	for _, query := range charSetQueries {
		err := Exec(query, charSet)
		exitOnError(err, "Error occurred setting connection character set.")
	}
	for _, query := range collationQueries {
		err := Exec(query, collation)
		exitOnError(err, "Error occurred setting connection collation.")
	}
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
	exitOnError(err, "Error occurred starting transaction.")
	tx = _tx
}

// Commit a transaction.
func Commit() {
	if tx != nil && tx.IsValid() {
		err := tx.Commit()
		exitOnError(err, "Error occurred committing transaction.")
	}
}

// Rollback a transaction.
func Rollback() {
	if tx != nil && tx.IsValid() {
		err := tx.Rollback()
		exitOnError(err, "Error occurred rolling back transaction.")
	}
}

// Create a database.
func createDatabase(name string, charSet string, collation string) (error) {
	err := ExecUnsafe("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET %s COLLATE %s;", name, charSet, collation)
	return err
}

// Drop a database.
func dropDatabase(name string) (error) {
	err := ExecUnsafe("DROP DATABASE IF EXISTS `%s`;", name)
	return err
}

// Change the database to the one named in the name parameter.
func useDatabase(name string) (error) {
	err := ExecUnsafe("USE `%s`;", name)
	return err
}

// Assert the database can be used. If not throw a fatal error.
func assertUseDatabase(name string) {
	err := useDatabase(name)
	exitOnError(err, "Can not use '%s' database.", name)
}
