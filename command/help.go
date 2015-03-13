package command

import (
	"strings"
	"text/tabwriter"
	"text/template"

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

func HelpPrinter(a *cli.App) func(templ string, data interface{}) {
	commonPrinter := func(templ string, data interface{}) {
		w := tabwriter.NewWriter(a.Writer, 0, 8, 1, '\t', 0)
		t := template.Must(template.New("help").Parse(templ))
		err := t.Execute(w, data)
		if err != nil {
			panic(err)
		}
		w.Flush()
	}
	return func(templ string, data interface{}) {
		if cmd, ok := data.(cli.Command); ok {
			type Help struct {
				cli.Command
				AppName  string
				Synopsis string
			}
			var d Help
			d.Command = cmd
			d.AppName = a.Name
			if s, ok := commandSynopsisses[cmd.Name]; ok {
				d.Synopsis = s
			}
			var desc string
			for _, ln := range strings.Split(strings.Trim(d.Description, helpTrimSpaces), "\n") {
				desc += helpPrefixSpaces + strings.Trim(ln, helpTrimSpaces) + "\n"
			}
			d.Description = strings.TrimRight(desc, helpTrimSpaces)
			commonPrinter(templ, d)
		} else {
			commonPrinter(templ, data)
		}
	}

}
