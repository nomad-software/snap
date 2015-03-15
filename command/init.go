// Package.
package command

// Imports.
import "github.com/codegangsta/cli"
import "github.com/nomad-software/snap/action"
import "log"

// Command.
var Init = cli.Command{
	Name:        "init",
	Usage:       "<database>",
	Description:
`Initialise a database to be managed.

ARGUMENTS:
    database
        The name of the database to be managed by snap.

EXAMPLE:

    snap init my_database
`,

	Action: func(ctx *cli.Context) {

		args := ctx.Args()

		if len(args) > 0 {
			name := args.First()
			action.InitialiseDatabase(name)
			return
		}

		log.Println("No database name specified.")
		log.Fatalf("Run '%s help init' for more information.\n", ctx.App.Name)
	},
}
