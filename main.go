// Package.
package main

// Imports.
import "github.com/codegangsta/cli"
import "os"
import "snap/command"

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

	app.Run(os.Args)
}
