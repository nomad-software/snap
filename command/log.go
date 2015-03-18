package command

// Imports.
import "github.com/codegangsta/cli"
import "github.com/nomad-software/snap/action"
import "log"

// Command.
var Log = cli.Command{
	Name:        "log",
	Usage:       "<database>",
	Description:
`Display a log of all schema update commits.

ARGUMENTS:
    database
        The name of the managed database to list commits for.

EXAMPLE:

    snap log my_database
`,

	Action: func(ctx *cli.Context) {

		args := ctx.Args()

		if len(args) > 0 {
			action.ShowLog(args.First())
			return
		}

		log.Println("No database name specified.")
		log.Fatalf("Run '%s help log' for more information.\n", ctx.App.Name)
	},
}
