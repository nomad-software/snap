// Package.
package sanitise

// Imports.
import "regexp"
import "strings"

// This function sanitises SQL for use by this program. The reason for this is 
// to prevent SQL being executed that would have un-desirable side effects or 
// cause errors. For example, when dealing with SQL updates many people craft 
// SQL for use on the command line. Such command line clients support SQL 
// features that this program does not. These features are removed from any SQL 
// executed by this program.
//
// All carriage returns are removed from the passed SQL string to help parsing. 
// All SQL lines containing the following statements are removed:
//
// 1. CREATE DATABASE
// 2. CREATE SCHEMA
// 3. USE
// 4. DELIMITER
func SanitiseSql(sql string) (string) {
	sql = ConvertToUnixLineEndings(sql)
	sql = removeCreateDatabaseStatements(sql)
	sql = removeUseStatements(sql)
	sql = reverseDelimiterChanges(sql)
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
// containing a DELIMITER statement and substitutes all occurances of custom 
// delimiters for the default semi-colon.
func reverseDelimiterChanges(sql string) (string) {
	// Match naked or quote delimited delimiters.
	pattern := regexp.MustCompile("(?i)DELIMITER\\s+(?:[`'\"](.*)[`'\"]|(\\S+))")
	lines   := strings.Split(sql, "\n")
	output  := make([]string, 0)
	var delimiter string;
	for _, line := range lines {
		if pattern.MatchString(line) {
			matches := pattern.FindStringSubmatch(line)
			if matches[1] != "" {
				// This capture group contains matched quoted delimiters.
				delimiter = matches[1]
			} else {
				// This capture group contains matched naked delimiters.
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
