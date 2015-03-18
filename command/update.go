package command

// Imports.
import "github.com/codegangsta/cli"
import "github.com/nomad-software/snap/action"
import "log"
import "strconv"

// Command.
var Update = cli.Command{
	Name:        "update",
	ShortName:   "up",
	Usage:       "<database> [revision]",
	Description:
`Update the database to a particular revision.

ARGUMENTS:
    database
        The name of the managed database to be updated to a particular
        schema revision.

    revision (optional)
        The schema revision to update the specified database to. This
        will default to the latest schema revision if not specified.

EXAMPLE:

    snap update my_database 10
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
			action.UpdateSchemaToRevision(database, revision)
			return
		}

		log.Println("No database name specified.")
		log.Fatalf("Run '%s help update' for more information.\n", ctx.App.Name)
	},
}
