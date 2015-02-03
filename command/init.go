// Package.
package command

// Imports.
import "github.com/codegangsta/cli"

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
		println("Args:", ctx.Args().First())
	},
}
