package command

// Imports.
import "github.com/codegangsta/cli"

// Command.
var Help = cli.Command{
	Name:        "help",
	Usage:       "[command]",
	Description:
`Shows a list of commands or extended help for one command.

ARGUMENTS:
    command (optional)
        The command to view extended help for.

EXAMPLE:

    snap help commit
`,

	Action: func(ctx *cli.Context) {
		args := ctx.Args()
		if args.Present() {
			cli.ShowCommandHelp(ctx, args.First())
		} else {
			cli.ShowAppHelp(ctx)
		}
	},
}

func init() {

	// Custom help message.
	cli.AppHelpTemplate = ` ___ _ __   __ _ _ __
/ __| '_ \ / _' | '_ \  v{{.Version}}
\__ \ | | | (_| | |_) |
|___/_| |_|\__,_| .__/
                |_|
{{.Usage}}
by {{.Author}} <{{.Email}}>

USAGE:
{{.Name}} command <arguments...> [optional]

COMMANDS:
{{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{"\t"}}{{.Usage}}
{{end}}
`

	// Custom command message.
	cli.CommandHelpTemplate = `COMMAND:
{{.Name}}{{with .ShortName}}, {{.}}{{end}}

USAGE:
snap {{.Name}} {{.Usage}}

DESCRIPTION:
{{.Description}}
`
}
