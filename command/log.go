// Package.
package command

// Imports.
import "github.com/codegangsta/cli"

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
		println("Args:", ctx.Args().First())
	},
}
