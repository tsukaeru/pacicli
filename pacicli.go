package main

import (
	"os"
	"runtime"

	"github.com/codegangsta/cli"
	"github.com/tsukaeru/pacicli/command"
)

var (
	commitHash string
	buildDate  string
)

func main() {
	app := cli.NewApp()
	app.Name = "pacicli"
	app.Version = Version + "\nGo version: " + runtime.Version()
	if len(commitHash) > 0 {
		app.Version += "\nGit commit: " + commitHash
	}
	if len(buildDate) > 0 {
		app.Version += "\nBuild date: " + buildDate
	}
	app.Usage = "Command line interface for Parallels Cloud Infrastructure"
	app.Commands = command.Commands
	app.EnableBashCompletion = true

	cli.CommandHelpTemplate = command.CommandHelpTemplate
	cli.HelpPrinter = command.HelpPrinter(app)
	app.Run(os.Args)
}
