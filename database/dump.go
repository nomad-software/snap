// Package.
package database

// Imports.
import "fmt"
import "strings"
import _ "github.com/go-sql-driver/mysql"

// Dump the entire schema of a database in SQL format to a string.
func DumpDatabase(databaseName string) (string) {
	AssertDatabaseExists(databaseName)
	sqlFragments := []string{
		createDatabaseDump(databaseName),
		createTablesDump(databaseName),
	}
	return strings.Join(sqlFragments, "\n\n");
}

// Create the tables dump string for an entire database.
func createTablesDump(databaseName string) (string) {
	AssertUseDatabase(databaseName)
	tables := getTablesForDatabase(databaseName)
	sqlFragments := make([]string, 0)
	for _, table := range tables {
		sqlFragments = append(sqlFragments, createTableDump(table))
	}
	return strings.Join(sqlFragments, "\n\n");
}

// Creat the table dump string for one table.
func createTableDump(tableName string) (string) {
	// TODO: Should test for existence of table!
	rows, err := db.Query(fmt.Sprintf("SHOW CREATE TABLE %s;", tableName))
	ExitOnError(err, fmt.Sprintf("Can not read table creation information for table '%s'.", tableName))
	var name string
	var sqlFragment string
	for rows.Next() {
		err = rows.Scan(&name, &sqlFragment)
		ExitOnError(err, fmt.Sprintf("Can not read creation sql for table '%s'.", tableName))
	}
	return sqlFragment
}

// Get all the tables inside the passed database.
func getTablesForDatabase(databaseName string) ([]string) {
	AssertUseDatabase(databaseName)
	rows, err := db.Query("SHOW TABLES;")
	ExitOnError(err, fmt.Sprintf("Can not access table information for database '%s'.", databaseName))
	var table string
	var tables = make([]string, 0)
	for rows.Next() {
		err = rows.Scan(&table)
		ExitOnError(err, fmt.Sprintf("Can not read table names for database '%s'.", databaseName))
		tables = append(tables, table)
	}
	return tables
}

// Create the database dump string.
func createDatabaseDump(databaseName string) string {
	query :=`SELECT
		SCHEMA_NAME,
		DEFAULT_CHARACTER_SET_NAME,
		DEFAULT_COLLATION_NAME
		FROM information_schema.SCHEMATA
		WHERE SCHEMA_NAME = ?
		LIMIT 1;`
	rows, err := db.Query(query, databaseName)
	ExitOnError(err, fmt.Sprintf("Can not access schema information for database '%s'.", databaseName))
	var name string
	var charSet string
	var collation string
	for rows.Next() {
		err = rows.Scan(&name, &charSet, &collation)
		ExitOnError(err, fmt.Sprintf("Can not read schema information for database '%s'.", databaseName))
	}
	return fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET %s COLLATE %s;", name, charSet, collation)
}
