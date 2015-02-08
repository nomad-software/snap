// Package.
package database

// Imports.
import "log"
import "os"
import _ "github.com/go-sql-driver/mysql"

// Constants.
const (
	configDatabaseName = "snap_config"
)

// Check if the snap config database exists.
func ConfigDatabaseExists() (bool) {
	return DatabaseExists(configDatabaseName)
}

// Switch to using the config database.
func UseConfigDatabase() {
	err := UseDatabase(configDatabaseName)
	ExitOnError(err, "Can not use config database.")
}

// Create the snap config database and all associated tables.
func CreateConfigDatabase() {
	success := true;
	success = success && createConfigDatabase()
	success = success && createConfigDatabasesTable()
	success = success && createConfigRevisionsTable()
	if !success {
		os.Exit(1)
	}
	log.Println("New snap config database created successfully")
}

// Create the config database.
func createConfigDatabase() (bool) {
	err := CreateDatabase(configDatabaseName)
	return WasSuccessful(err)
}

// Create the config database's database table.
func createConfigDatabasesTable() (bool) {
	sql := `
CREATE TABLE IF NOT EXISTS ` + configDatabaseName + `.initialisedDatabases (
  id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(64) NOT NULL,
  dateInitialized TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE INDEX uniqueDatabaseName (name ASC))
ENGINE = InnoDB;
`
	_, err := db.Exec(sql)
	return WasSuccessful(err)
}

// Create the config database's revisions table.
func createConfigRevisionsTable() (bool) {
	sql := `
CREATE TABLE IF NOT EXISTS ` + configDatabaseName + `.revisions (
  id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  databaseId INT UNSIGNED NOT NULL,
  revision INT UNSIGNED NOT NULL,
  upSql TEXT NOT NULL,
  downSql TEXT NOT NULL,
  fullSql TEXT NOT NULL COMMENT 'SQL snapshot after applying update SQL.',
  comment VARCHAR(255) NOT NULL,
  dateApplied TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  INDEX databaseId_FK (databaseId ASC),
  UNIQUE INDEX uniqueDatabaseIdAndRevision (databaseId ASC, revision ASC),
  CONSTRAINT databaseIdForeignKey
    FOREIGN KEY (databaseId)
    REFERENCES ` + configDatabaseName + `.initialisedDatabases (id)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;
`
	_, err := db.Exec(sql)
	return WasSuccessful(err)
}
