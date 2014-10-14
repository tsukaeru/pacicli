package command

import (
	"strings"

	"github.com/codegangsta/cli"
)

const CommandHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.AppName}} {{.Name}}{{if .Synopsis}} {{.Synopsis}}{{else}}{{if .Flags}} [command options]{{end}} [arguments...]{{end}}{{if .Description}}

DESCRIPTION:
{{.Description}}{{end}}{{if .Flags}}

OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{ end }}
`

const (
	helpTrimSpaces   = "\n\t "
	helpPrefixSpaces = "   "
)

func HelpPrinter(appName string) func(templ string, data interface{}) {
	origPrinter := cli.HelpPrinter
	return func(templ string, data interface{}) {
		if cmd, ok := data.(cli.Command); ok {
			type Help struct {
				cli.Command
				AppName  string
				Synopsis string
			}
			var d Help
			d.Command = cmd
			d.AppName = appName
			if s, ok := commandSynopsisses[cmd.Name]; ok {
				d.Synopsis = s
			}
			var desc string
			for _, ln := range strings.Split(strings.Trim(d.Description, helpTrimSpaces), "\n") {
				desc += helpPrefixSpaces + strings.Trim(ln, helpTrimSpaces) + "\n"
			}
			d.Description = strings.TrimRight(desc, helpTrimSpaces)
			origPrinter(templ, d)
		} else {
			origPrinter(templ, data)
		}
	}

}
