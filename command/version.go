// Package.
package command

// Imports.
import "github.com/codegangsta/cli"

// Command.
var Version = cli.Command{
	Name:        "version",
	Usage:       "",
	Description:
`Displays version information about the executable.

EXAMPLE:

    snap version
`,

	Action: func(ctx *cli.Context) {
		println(ctx.App.Name, ctx.App.Version)
	},
}
