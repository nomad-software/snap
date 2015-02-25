// Package.
package command

// Imports.
import "github.com/codegangsta/cli"
import "github.com/nomad-software/snap/action"
import "log"
import "strconv"

// Command.
var Show = cli.Command{
	Name:        "show",
	Usage:       "<database> [revision]",
	Description:
`Show the update SQL for a particular revision. This is the SQL that appeared 
in the original snap file within the UP section.

ARGUMENTS:
    database
        The name of the managed database.

    revision (optional)
        The schema revision to show the update SQL of. This will default
        to the latest schema revision if not specified.

EXAMPLE:

    snap show my_database 10
`,

	Action: func(ctx *cli.Context) {

		args := ctx.Args()

		if len(args) > 0 {
			database := args.Get(0)
			// Ignore the error when getting the second argument because if the 
			// argument can not be parsed to a uint then (along with the error) 
			// zero is returned, which is what we want because we can use it as 
			// an empty value.
			revision, _ := strconv.ParseUint(args.Get(1), 10, 64)
			action.ShowUpdateSql(database, revision)
			return
		}

		log.Println("No database name specified")
		log.Fatalf("Run '%s help show' for more information\n", ctx.App.Name)
	},
}
