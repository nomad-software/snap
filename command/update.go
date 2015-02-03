// Package.
package command

// Imports.
import "github.com/codegangsta/cli"

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
		println("Args:", ctx.Args().First())
	},
}
