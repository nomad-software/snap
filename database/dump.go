// Package.
package database

// Imports.
import "fmt"
import "strings"

// Dump the entire schema of a database in SQL format to a string.
func DumpDatabase(databaseName string) (string) {
	AssertUseDatabase(databaseName)
	output := []string{
		exportDatabase(databaseName),
		exportTables(databaseName),
		exportFunctions(databaseName),
		exportProcedures(databaseName),
		exportTriggers(databaseName),
	}
	// Filter out empty output.
	sqlFragments := make([]string, 0)
	for _, sqlFragment := range output {
		if sqlFragment != "" {
			sqlFragments = append(sqlFragments, sqlFragment)
		}
	}
	return strings.Join(sqlFragments, "\n\n");
}

// Generate a comment to separate the sections.
func generateCommentHeading(heading string) (string) {
	line := "-- +----------------------------------------------------------------------------"
	return fmt.Sprintf("%s\n-- | %s\n%s", line, heading, line)
}

// Prepend a comment header to the passed slice of SQL fragments. If the slice 
// is empty, just return the empty slice.
func prependHeaderFragment(heading string, sqlFragments []string) ([]string) {
	if len(sqlFragments) > 0 {
		return append([]string{generateCommentHeading(heading)}, sqlFragments...)
	}
	return sqlFragments
}

// Wrap delimiter sensitive SQL fragments with safe delimiters. If the passed 
// slice is empty, just return it.
func wrapFragmentsWithSafeDelimiters(sqlFragments []string) ([]string) {
	if len(sqlFragments) > 0 {
		sqlFragments = append([]string{"DELIMITER $$"}, sqlFragments...)
		sqlFragments = append(sqlFragments, "DELIMITER ;")
	}
	return sqlFragments
}

// Export the database SQL. This function assumes the database exists and is 
// being used.
func exportDatabase(databaseName string) string {
	query :=`SELECT
		DEFAULT_CHARACTER_SET_NAME,
		DEFAULT_COLLATION_NAME
		FROM information_schema.SCHEMATA
		WHERE SCHEMA_NAME = ?
		LIMIT 1;`
	row, err := QueryRow(query, databaseName)
	ExitOnError(err, fmt.Sprintf("Can not access schema information for database '%s'.", databaseName))
	sqlFragments := make([]string, 0)
	sqlFragments = append(sqlFragments, fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET %s COLLATE %s;", databaseName, row.Str(0), row.Str(1)))
	sqlFragments = append(sqlFragments, fmt.Sprintf("USE DATABASE `%s`;", databaseName))
	sqlFragments = prependHeaderFragment("Database", sqlFragments)
	return strings.Join(sqlFragments, "\n\n");
}

// Export table SQL string for an entire database. This function assumes the 
// database exists and is being used.
func exportTables(databaseName string) (string) {
	tables := getAllTableNames(databaseName)
	sqlFragments := make([]string, 0)
	for _, table := range tables {
		sqlFragments = append(sqlFragments, exportTable(table))
	}
	sqlFragments = prependHeaderFragment("Tables", sqlFragments)
	return strings.Join(sqlFragments, "\n\n");
}

// Retrieve all the table names from the passed database. This function assumes 
// the database exists and is being used.
func getAllTableNames(databaseName string) ([]string) {
	rows, err := Query("SHOW TABLES;")
	ExitOnError(err, fmt.Sprintf("Can not access table information for database '%s'.", databaseName))
	var tables = make([]string, 0)
	for _, row := range rows {
		tables = append(tables, row.Str(0))
	}
	return tables
}

// Export the SQL for one table. This function assumes the table exists.
func exportTable(tableName string) (string) {
	row, err := QueryRowUnsafe("SHOW CREATE TABLE `%s`;", tableName)
	ExitOnError(err, fmt.Sprintf("Can not read creation information for table '%s'.", tableName))
	// The ending semi-colon is always missing when retrieving an SQL fragment 
	// like this.
	return row.Str(1) + ";"
}

// Export function SQL string for an entire database. This function assumes the 
// database exists and is being used.
func exportFunctions(databaseName string) (string) {
	functions := getAllFunctionNames(databaseName)
	sqlFragments := make([]string, 0)
	for _, function := range functions {
		sqlFragments = append(sqlFragments, exportFunction(function))
	}
	sqlFragments = wrapFragmentsWithSafeDelimiters(sqlFragments)
	sqlFragments = prependHeaderFragment("Functions", sqlFragments)
	return strings.Join(sqlFragments, "\n\n");
}

// Retrieve all the function names from the passed database. This function 
// assumes the database exists and is being used.
func getAllFunctionNames(databaseName string) ([]string) {
	rows, err := Query("SHOW FUNCTION STATUS WHERE Db = ?;", databaseName)
	ExitOnError(err, fmt.Sprintf("Can not access function information for database '%s'.", databaseName))
	var functions = make([]string, 0)
	for _, row := range rows {
		functions = append(functions, row.Str(1))
	}
	return functions
}

// Export the SQL for one function. This function assumes the function exists.
func exportFunction(functionName string) (string) {
	row, err := QueryRowUnsafe("SHOW CREATE FUNCTION `%s`;", functionName)
	ExitOnError(err, fmt.Sprintf("Can not read creation information for function '%s'.", functionName))
	// The ending safe delimiter is always missing when retrieving an SQL 
	// fragment like this.
	return row.Str(2) + "$$"
}

// Export procedure SQL string for an entire database. This function assumes 
// the database exists and is being used.
func exportProcedures(databaseName string) (string) {
	procedures := getAllProcedureNames(databaseName)
	sqlFragments := make([]string, 0)
	for _, procedure := range procedures {
		sqlFragments = append(sqlFragments, exportProcedure(procedure))
	}
	sqlFragments = wrapFragmentsWithSafeDelimiters(sqlFragments)
	sqlFragments = prependHeaderFragment("Procedures", sqlFragments)
	return strings.Join(sqlFragments, "\n\n");
}

// Retrieve all the procedure names from the passed database. This function 
// assumes the database exists and is being used.
func getAllProcedureNames(databaseName string) ([]string) {
	rows, err := Query("SHOW PROCEDURE STATUS WHERE Db = ?;", databaseName)
	ExitOnError(err, fmt.Sprintf("Can not access procedure information for database '%s'.", databaseName))
	var procedures = make([]string, 0)
	for _, row := range rows {
		procedures = append(procedures, row.Str(1))
	}
	return procedures
}

// Export the SQL for one procedure. This function assumes the procedure 
// exists.
func exportProcedure(procedureName string) (string) {
	row, err := QueryRowUnsafe("SHOW CREATE PROCEDURE `%s`;", procedureName)
	ExitOnError(err, fmt.Sprintf("Can not read creation information for procedure '%s'.", procedureName))
	// The ending safe delimiter is always missing when retrieving an SQL 
	// fragment like this.
	return row.Str(2) + "$$"
}

// Export trigger SQL string for an entire database. This function assumes the 
// database exists and is being used.
func exportTriggers(databaseName string) (string) {
	triggers := getAllTriggerNames(databaseName)
	sqlFragments := make([]string, 0)
	for _, trigger := range triggers {
		sqlFragments = append(sqlFragments, exportTrigger(trigger))
	}
	sqlFragments = wrapFragmentsWithSafeDelimiters(sqlFragments)
	sqlFragments = prependHeaderFragment("Triggers", sqlFragments)
	return strings.Join(sqlFragments, "\n\n");
}

// Retrieve all the trigger names from the passed database. This function 
// assumes the database exists and is being used.
func getAllTriggerNames(databaseName string) ([]string) {
	rows, err := QueryUnsafe("SHOW TRIGGERS FROM `%s`;", databaseName)
	ExitOnError(err, fmt.Sprintf("Can not access trigger information for database '%s'.", databaseName))
	var triggers = make([]string, 0)
	for _, row := range rows {
		triggers = append(triggers, row.Str(0))
	}
	return triggers
}

// Export the SQL for one trigger. This function assumes the trigger exists.
func exportTrigger(triggerName string) (string) {
	row, err := QueryRowUnsafe("SHOW CREATE TRIGGER `%s`;", triggerName)
	ExitOnError(err, fmt.Sprintf("Can not read creation information for trigger '%s'.", triggerName))
	// The ending safe delimiter is always missing when retrieving an SQL 
	// fragment like this.
	return row.Str(2) + "$$"
}
