// Package.
package command

// Imports.
import "github.com/codegangsta/cli"

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
		println("Args:", ctx.Args().First())
	},
}
