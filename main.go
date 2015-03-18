// Snap is a proof of concept tool to start exploring version control for 
// database schemas. Usually when maintaining and updating database schemas 
// over multiple environments things start to get confusing very quickly. Snap 
// is a tool inspired by git allowing you to manage and interrogate snap 
// managed databases.
package main

// Imports.
import "github.com/codegangsta/cli"
import "github.com/nomad-software/snap/command"
import "github.com/nomad-software/snap/config"
import "github.com/nomad-software/snap/database"
import "os"

// Main entry point.
func main() {

	app := cli.NewApp()

	app.Name        = "snap"
	app.Version     = "0.1a"
	app.Usage       = "Version control for MySql database schemas"
	app.Author      = "Gary Willoughby"
	app.Email       = "snap@nomad.so"
	app.HideHelp    = true
	app.HideVersion = true

	app.Commands = []cli.Command{
		command.Commit,
		command.Copy,
		command.Diff,
		command.Dump,
		command.Help,
		command.Init,
		command.List,
		command.Log,
		command.Show,
		command.Update,
		command.Version,
	}

	app.Action = func(ctx *cli.Context) {
		cli.ShowAppHelp(ctx)
	}

	database.Open(config.GetConfig())
	defer database.Close()

	app.Run(os.Args)
}
