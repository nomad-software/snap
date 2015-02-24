// Package.
package command

// Imports.
import "fmt"
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
		if len(ctx.Args()) > 0 {
			name := ctx.Args().First()
			action.ShowLog(name)
			return
		}
		log.Println("No database name specified")
		log.Fatalln(fmt.Sprintf("Run '%s help log' for more information", ctx.App.Name))
	},
}
