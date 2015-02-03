// Package.
package command

// Imports.
import "github.com/codegangsta/cli"

// Command.
var List = cli.Command{
	Name:        "list",
	ShortName:   "ls",
	Usage:       "",
	Description:
`List the name of all databases snap is currently managing.

EXAMPLE:

    snap list
`,

	Action: func(ctx *cli.Context) {
		println("Args:", ctx.Args().First())
	},
}
