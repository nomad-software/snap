// Package.
package config

// Imports.
import "encoding/json"
import "fmt"
import "github.com/mitchellh/go-homedir"
import "io/ioutil"
import "log"

// This struct holds the database configuration details.
type Database struct {
	User string
	Password string
	Protocol string
	Host string
	Port string
}

// Format the Database struct into a valid DSN (data source name) string.
func (this Database) String() (string) {
	return fmt.Sprintf("%s:%s@%s(%s:%s)/",
		this.User,
		this.Password,
		this.Protocol,
		this.Host,
		this.Port)
}

// Return a new Config struct initialised with default values.
func NewDatabase() Database {
	return Database{
		Protocol: "tcp",
		Host: "127.0.0.1",
		Port: "3306",
	}
}

// Get the full path of the config file.
func getConfigFilePath() (string) {
	path, err := homedir.Expand("~/.snap")
	if err != nil {
		log.Fatalln("Can not determine user's home directory.")
	}
	return path
}

// Read the config file.
func readConfigFile(path string) ([]byte) {
	
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(`The snap config file is missing: (~/.snap). This file should be in the following Json format:

{
   "user":"foo",
   "password": "bar",
   "protocol": "tcp",
   "host": "localhost",
   "port": "3306"
}

The protocol, host and port fields are optional and default to the values shown above.
`)
	}
	return file
}

// Parse the config file into a Database struct.
func ParseConfigFile() (Database) {

	path     := getConfigFilePath()
	contents := readConfigFile(path);
	config   := NewDatabase();

	err := json.Unmarshal(contents, &config)
	if err != nil {
		log.Fatalln("Can not un-marshal json into config struct.")
	}

	return config
}
