package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/tsukaeru/pacicli/command"
)

func main() {
	app := cli.NewApp()
	app.Name = "pacicli"
	app.Version = Version
	app.Usage = "Command line interface for Parallels Cloud Infrastructure"
	app.Commands = command.Commands
	app.EnableBashCompletion = true

	cli.CommandHelpTemplate = command.CommandHelpTemplate
	cli.HelpPrinter = command.HelpPrinter(app.Name)
	app.Run(os.Args)
}
