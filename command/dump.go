// Package.
package command

// Imports.
import "github.com/codegangsta/cli"
import "github.com/nomad-software/snap/action"
import "log"
import "strconv"

// Command.
var Dump = cli.Command{
	Name:        "dump",
	Usage:       "<database> [revision]",
	Description:
`Dump a managed database schema to stdout.

ARGUMENTS:
    database
        The name of the managed database to dump the schema from.

    revision (optional)
        The schema revision from which to dump the schema. This will
        default to the latest schema revision if not specified.

EXAMPLE:

    snap dump my_database 10
`,

	Action: func(ctx *cli.Context) {
		args := ctx.Args()

		if len(args) > 0 {
			database := args.Get(0)
			// Ignore the error when getting the second argument because if the 
			// argument can not be parsed to a uint64 then (along with the 
			// error) zero is returned, which is what we want because we can 
			// use it as an empty value.
			revision, _ := strconv.ParseUint(args.Get(1), 10, 64)
			action.ShowFullSql(database, revision)
			return
		}

		log.Println("No database name specified")
		log.Fatalf("Run '%s help dump' for more information\n", ctx.App.Name)
	},
}
