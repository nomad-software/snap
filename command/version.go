// Package.
package command

// Imports.
import "fmt"
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
		fmt.Println(ctx.App.Name, ctx.App.Version)
	},
}
