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

// Create the snap config database and all associated tables.
func CreateConfigDatabase() {
	success := true;
	success = success && createConfigDatabase()
	success = success && createConfigDatabaseTable()
	if !success {
		os.Exit(1)
	}
	log.Println("New snap config database created successfully")
}

// Create the config database.
func createConfigDatabase() (bool) {
	err = CreateDatabase("snap_config")
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// Create the config database's database table.
func createConfigDatabaseTable() (bool) {
	sql := `CREATE TABLE IF NOT EXISTS snap_config.initialisedDatabases (
  id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(64) NOT NULL,
  initialized TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE INDEX databaseName_UNIQUE (name ASC))
ENGINE = InnoDB;`
	_, err = db.Query(sql)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
