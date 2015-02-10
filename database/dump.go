// Package.
package database

// Imports.
import "fmt"
import "strings"
import _ "github.com/go-sql-driver/mysql"

// Dump the entire schema of a database in SQL format to a string.
func DumpDatabase(databaseName string) (string) {
	AssertUseDatabase(databaseName)
	sqlFragments := []string{
		generateDatabaseDump(databaseName),
		generateTablesDump(databaseName),
		generateFunctionsDump(databaseName),
		generateProceduresDump(databaseName),
		generateTriggersDump(databaseName),
	}
	return strings.Join(sqlFragments, "\n\n");
}

// Generate a create database string for the passed database.
// This function assumes the database exists and is being used.
func generateDatabaseDump(databaseName string) string {
	query :=`SELECT
		SCHEMA_NAME,
		DEFAULT_CHARACTER_SET_NAME,
		DEFAULT_COLLATION_NAME
		FROM information_schema.SCHEMATA
		WHERE SCHEMA_NAME = ?
		LIMIT 1;`
	rows, err := db.Query(query, databaseName)
	ExitOnError(err, fmt.Sprintf("Can not access schema information for database '%s'.", databaseName))
	var name, charSet, collation []byte
	for rows.Next() {
		err = rows.Scan(&name, &charSet, &collation)
		ExitOnError(err, fmt.Sprintf("Can not read schema information for database '%s'.", databaseName))
	}
	return fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET %s COLLATE %s;", string(name), string(charSet), string(collation))
}

// Create the tables dump string for an entire database.
// This function assumes the database exists and is being used.
func generateTablesDump(databaseName string) (string) {
	tables := getAllTables(databaseName)
	sqlFragments := make([]string, 0)
	for _, table := range tables {
		sqlFragments = append(sqlFragments, generateSingleTableDump(table))
	}
	return strings.Join(sqlFragments, "\n\n");
}

// Get all the tables inside the passed database.
// This function assumes the database exists and is being used.
func getAllTables(databaseName string) ([]string) {
	rows, err := db.Query("SHOW TABLES;")
	ExitOnError(err, fmt.Sprintf("Can not access table information for database '%s'.", databaseName))
	var table []byte
	var tables = make([]string, 0)
	for rows.Next() {
		err = rows.Scan(&table)
		ExitOnError(err, fmt.Sprintf("Can not read table names for database '%s'.", databaseName))
		tables = append(tables, string(table))
	}
	return tables
}

// Create the table dump string for one table.
// This function assumes the table exists.
func generateSingleTableDump(tableName string) (string) {
	rows, err := db.Query(fmt.Sprintf("SHOW CREATE TABLE %s;", tableName))
	ExitOnError(err, fmt.Sprintf("Can not read creation information for table '%s'.", tableName))
	var name, sqlFragment []byte
	for rows.Next() {
		err = rows.Scan(&name, &sqlFragment)
		ExitOnError(err, fmt.Sprintf("Can not read creation sql for table '%s'.", tableName))
	}
	return string(sqlFragment)
}

// Create the functions dump string for an entire database.
// This function assumes the database exists and is being used.
func generateFunctionsDump(databaseName string) (string) {
	functions := getAllFunctions(databaseName)
	sqlFragments := make([]string, 0)
	for _, function := range functions {
		sqlFragments = append(sqlFragments, generateSingleFunctionDump(function))
	}
	return strings.Join(sqlFragments, "\n\n");
}

// Get all the functions inside the passed database.
// This function assumes the database exists and is being used.
func getAllFunctions(databaseName string) ([]string) {
	rows, err := db.Query("SHOW FUNCTION STATUS WHERE Db = ?;", databaseName)
	ExitOnError(err, fmt.Sprintf("Can not access function information for database '%s'.", databaseName))
	var database, name, type_, definer, modified, created, security, comment, charSet, collationCon, collationDb []byte
	var functions = make([]string, 0)
	for rows.Next() {
		err = rows.Scan(&database, &name, &type_, &definer, &modified, &created, &security, &comment, &charSet, &collationCon, &collationDb)
		ExitOnError(err, fmt.Sprintf("Can not read function names for database '%s'.", databaseName))
		functions = append(functions, string(name))
	}
	return functions
}

// Create the function dump string for one function.
// This function assumes the function exists.
func generateSingleFunctionDump(functionName string) (string) {
	rows, err := db.Query(fmt.Sprintf("SHOW CREATE FUNCTION %s;", functionName))
	ExitOnError(err, fmt.Sprintf("Can not read creation information for function '%s'.", functionName))
	var name, mode, sqlFragment, charSet, collationCon, collationDb []byte
	for rows.Next() {
		err = rows.Scan(&name, &mode, &sqlFragment, &charSet, &collationCon, &collationDb)
		ExitOnError(err, fmt.Sprintf("Can not read creation sql for function '%s'.", functionName))
	}
	return string(sqlFragment)
}

// Create the procedures dump string for an entire database.
// This function assumes the database exists and is being used.
func generateProceduresDump(databaseName string) (string) {
	procedures := getAllProcedures(databaseName)
	sqlFragments := make([]string, 0)
	for _, procedure := range procedures {
		sqlFragments = append(sqlFragments, generateSingleProcedureDump(procedure))
	}
	return strings.Join(sqlFragments, "\n\n");
}

// Get all the procedures inside the passed database.
// This function assumes the database exists and is being used.
func getAllProcedures(databaseName string) ([]string) {
	rows, err := db.Query("SHOW PROCEDURE STATUS WHERE Db = ?;", databaseName)
	ExitOnError(err, fmt.Sprintf("Can not access procedure information for database '%s'.", databaseName))
	var database, name, type_, definer, modified, created, security, comment, charSet, collationCon, collationDb []byte
	var procedures = make([]string, 0)
	for rows.Next() {
		err = rows.Scan(&database, &name, &type_, &definer, &modified, &created, &security, &comment, &charSet, &collationCon, &collationDb)
		ExitOnError(err, fmt.Sprintf("Can not read procedure names for database '%s'.", databaseName))
		procedures = append(procedures, string(name))
	}
	return procedures
}

// Create the procedure dump string for one procedure.
// This function assumes the procedure exists.
func generateSingleProcedureDump(procedureName string) (string) {
	rows, err := db.Query(fmt.Sprintf("SHOW CREATE PROCEDURE %s;", procedureName))
	ExitOnError(err, fmt.Sprintf("Can not read creation information for procedure '%s'.", procedureName))
	var name, mode, sqlFragment, charSet, collationCon, collationDb []byte
	for rows.Next() {
		err = rows.Scan(&name, &mode, &sqlFragment, &charSet, &collationCon, &collationDb)
		ExitOnError(err, fmt.Sprintf("Can not read creation sql for procedure '%s'.", procedureName))
	}
	return string(sqlFragment)
}

// Create the triggers dump string for an entire database.
// This function assumes the database exists and is being used.
func generateTriggersDump(databaseName string) (string) {
	triggers := getAllTriggers(databaseName)
	sqlFragments := make([]string, 0)
	for _, trigger := range triggers {
		sqlFragments = append(sqlFragments, generateSingleTriggerDump(trigger))
	}
	return strings.Join(sqlFragments, "\n\n");
}

// Get all the triggers inside the passed database.
// This function assumes the database exists and is being used.
func getAllTriggers(databaseName string) ([]string) {
	rows, err := db.Query(fmt.Sprintf("SHOW TRIGGERS FROM %s;", databaseName))
	ExitOnError(err, fmt.Sprintf("Can not access trigger information for database '%s'.", databaseName))
	var name, event, table, statement, timing, created, mode, definer, charSet, collationCon, collationDb []byte
	var triggers = make([]string, 0)
	for rows.Next() {
		err = rows.Scan(&name, &event, &table, &statement, &timing, &created, &mode, &definer, &charSet, &collationCon, &collationDb)
		ExitOnError(err, fmt.Sprintf("Can not read trigger names for database '%s'.", databaseName))
		triggers = append(triggers, string(name))
	}
	return triggers
}

// Create the trigger dump string for one trigger.
// This function assumes the trigger exists.
func generateSingleTriggerDump(triggerName string) (string) {
	rows, err := db.Query(fmt.Sprintf("SHOW CREATE TRIGGER %s;", triggerName))
	ExitOnError(err, fmt.Sprintf("Can not read creation information for trigger '%s'.", triggerName))
	var name, mode, sqlFragment, charSet, collationCon, collationDb []byte
	for rows.Next() {
		err = rows.Scan(&name, &mode, &sqlFragment, &charSet, &collationCon, &collationDb)
		ExitOnError(err, fmt.Sprintf("Can not read creation sql for trigger '%s'.", triggerName))
	}
	return string(sqlFragment)
}
