// Package.
package command

// Imports.
import "github.com/codegangsta/cli"

// Command.
var Commit = cli.Command{
	Name:        "commit",
	ShortName:   "ci",
	Usage:       "<database> <snapfile> <message>",
	Description:
`Commit a new schema revision to a managed database. A schema revision is 
defined within a snap file which follows the format described below. This file 
is then applied to the database.

ARGUMENTS:
    database
        The managed database to apply the snap file to.

    snapfile
        The file holding the changes to be applied.

    message
        The message to store against the commit.

SNAPFILE:
A snap file is a simple text file holding SQL statements to be applied to the 
database. Two textual delimeters are required in the file. The first (marked 
using UP:) is used to mark the start of SQL making changes to the database 
schema. The second (marked using DOWN:) is used to mark the start of SQL 
reversing the changes made in the first section. Here is a sample file:

    UP:
    CREATE TABLE IF NOT EXISTS foo (
      bar TINYINT UNSIGNED NOT NULL,
      baz VARCHAR(32) NOT NULL,
      PRIMARY KEY (id) )
    ENGINE = InnoDB;

    DOWN:
    DROP TABLE IF EXISTS foo;

Once the commit is successful the snap file can be discarded as it is saved to 
the snap configuration database.

EXAMPLE:

    snap my_database changes.txt "Added table foo."

WARNING:
It is the developer's responsibility to ensure the SQL in the DOWN 
section correctly reverses the changes made in the UP section. Failing to 
ensure this could corrupt the database.
	`,

	Action: func(ctx *cli.Context) {
		println("Args:", ctx.Args().First())
	},
}
