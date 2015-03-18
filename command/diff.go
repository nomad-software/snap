package command

// Imports.
import "github.com/codegangsta/cli"
import "github.com/nomad-software/snap/action"
import "log"

// Command.
var Diff = cli.Command{
	Name:       "diff",
	Usage:      "<database> <from-revision>[..<to-revision>]",
	Description:
`Show an SQL diff between two schema revisions. The diff will be in unified 
format and be written to stdout.

ARGUMENTS:
    database
        The name of the managed database to generate the schema diff for.

    from-revision
        A schema revision to generate the diff from. This is usually the
        start point in history where you would like the diff to be
        calculated from.

    to-revision (optional)
        The schema revision to generate the diff to. This marks the finish
        point in history where you would like the diff to be calculated to.
        This will default to the latest schema revision if not specified.

EXAMPLE:

	snap diff my_database 10..12
`,

	Action: func(ctx *cli.Context) {
		args := ctx.Args()

		if len(args) > 1 {
			database       := args.Get(0)
			revisionString := args.Get(1)
			action.Diff(database, revisionString)
			return
		}

		log.Println("Both database and a revision must be specified.")
		log.Fatalf("Run '%s help diff' for more information.\n", ctx.App.Name)
	},
}
