// Package.
package command

// Imports.
import "github.com/codegangsta/cli"

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
		println("Args:", ctx.Args().First())
	},
}
