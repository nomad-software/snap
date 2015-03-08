// Package.
package database

// Imports.
import "regexp"
import "strings"

// Sanitises SQL for use by this program. The reason for sanitisation is to 
// prevent SQL being executed by this program that would have un-desirable side 
// effects. For example, when dealing with SQL updates many people use the 
// command line. Such command line clients support features that this program 
// does not. These features are removed from any SQL executed by this program.
//
// All carriage returns are removed. All SQL lines containing the following 
// statements are removed:
//
// 1. CREATE DATABASE
// 2. CREATE SCHEMA
// 3. USE
// 4. DELIMITER
func sanitise(sql string) (string) {
	sql = convertToUnixLineEndings(sql)
	sql = removeCreateDatabaseStatements(sql)
	sql = removeUseStatements(sql)
	sql = reverseDelimiterChanges(sql)
	return sql
}

// Convert the SQL string line endings to unix format.
func convertToUnixLineEndings(sql string) (string) {
	sql = strings.Replace(sql, "\r\n", "\n", -1);
	sql = strings.Replace(sql, "\r", "\n", -1);
	return sql
}

// Remove any CREATE DATABASE or CREATE SCHEMA statements in the passed SQL.
func removeCreateDatabaseStatements(sql string) (string) {
	pattern := regexp.MustCompile("(?i)(CREATE DATABASE\\s+|CREATE SCHEMA\\s+)")
	lines   := strings.Split(sql, "\n")
	output  := make([]string, 0)
	for _, line := range lines {
		if !pattern.MatchString(line) {
			output = append(output, line)
		}
	}
	return strings.Join(output, "\n")
}

// Remove any USE statements in the passed SQL.
func removeUseStatements(sql string) (string) {
	pattern := regexp.MustCompile("(?i)USE\\s+")
	lines   := strings.Split(sql, "\n")
	output  := make([]string, 0)
	for _, line := range lines {
		if !pattern.MatchString(line) {
			output = append(output, line)
		}
	}
	return strings.Join(output, "\n")
}

// Reverse any delimiter changes in the passed SQL. This removes all lines 
// containing a DELIMITER statement and replaces all occurances of that custom 
// delimiter to the default semi-colon.
func reverseDelimiterChanges(sql string) (string) {
	pattern := regexp.MustCompile("(?i)DELIMITER\\s+(?:[`'\"](.*)[`'\"]|(\\S+))")
	lines   := strings.Split(sql, "\n")
	output  := make([]string, 0)
	var delimiter string;
	for _, line := range lines {
		if pattern.MatchString(line) {
			matches := pattern.FindStringSubmatch(line)
			if matches[1] != "" {
				delimiter = matches[1]
			} else {
				delimiter = matches[2]
			}
		} else {
			if delimiter != "" {
				line = strings.Replace(line, delimiter, ";", -1)
			}
			output = append(output, line)
		}
	}
	return strings.Join(output, "\n")
}
