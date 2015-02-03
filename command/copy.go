// Package.
package command

// Imports.
import "github.com/codegangsta/cli"

// Command.
var Copy = cli.Command{
	Name:        "copy",
	ShortName:   "cp",
	Usage:       "<source-database> <destination-database> [revision]",
	Description:
`Copy a database schema to a new database.

ARGUMENTS:
    source-database
        The name of the managed database which holds the source schema.

    destination-database
        The new database which will be created with the schema of the source
        database. This new database will not be initially managed by snap.

    revision (optional)
        The schema revision of the source database to use for creating the
        new database. This will default to the latest schema revision if not
        specified.

EXAMPLE:

    snap copy my_database my_new_database 12
	`,

	Action: func(ctx *cli.Context) {
		println("Args:", ctx.Args().First())
	},
}
