// Package.
package command

// Imports.
import "github.com/codegangsta/cli"
import "github.com/nomad-software/snap/action"
import "log"
import "strconv"

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
		args := ctx.Args()

		if len(args) > 1 {
			source      := args.Get(0)
			destination := args.Get(1)
			// Ignore the error when getting the third argument because if the 
			// argument can not be parsed to a uint64 then (along with the 
			// error) zero is returned, which is what we want because we can 
			// use it as an empty value.
			revision, _ := strconv.ParseUint(args.Get(2), 10, 64)
			action.CopyDatabase(source, destination, revision)
			return
		}

		log.Println("Both source and destination databases must be specified")
		log.Fatalf("Run '%s help copy' for more information\n", ctx.App.Name)
	},
}
