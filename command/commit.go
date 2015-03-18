package command

// Imports.
import "github.com/codegangsta/cli"
import "github.com/nomad-software/snap/action"
import "log"

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
database. Two SQL comments are required in the file to act as delimiters. The 
first, marked '-- SNAP_UP', is used to mark the start of SQL making any 
required changes to the database schema. The second, marked '-- SNAP_DOWN', is 
used to mark the start of SQL reversing the changes made in the first section. 
Both comments must be on their own line. Here is a sample file:

    -- SNAP_UP
    CREATE TABLE IF NOT EXISTS foo (
      bar TINYINT UNSIGNED NOT NULL,
      baz VARCHAR(32) NOT NULL,
      PRIMARY KEY (bar) )
    ENGINE = InnoDB;

    -- SNAP_DOWN
    DROP TABLE IF EXISTS foo;

Once a commit is successful the snap file can be discarded as it is saved to 
the snap configuration database.

EXAMPLE:

    snap my_database changes.txt "Added table foo."
	`,

	Action: func(ctx *cli.Context) {
		args := ctx.Args()

		if len(args) > 2 {
			database := args.Get(0)
			fileName := args.Get(1)
			message  := args.Get(2)
			action.CommitFile(database, fileName, message)
			return
		}

		log.Println("Not enough arguments supplied.")
		log.Fatalf("Run '%s help commit' for more information.\n", ctx.App.Name)
	},
}
