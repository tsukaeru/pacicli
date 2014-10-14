package command

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/tatsushid/go-prettytable"
	"github.com/tsukaeru/pacicli/lib"
)

var commandApplicationList = cli.Command{
	Name:      "applist",
	ShortName: "appls",
	Usage:     "List application templates",
	Description: `
	This command obtains a list of the available application templates for
	Container.
`,
	Flags: append(CommonFlags, noHeaderFlag),
	Action: func(c *cli.Context) {
		action(c, doApplicationList)
	},
}

var commandApplicationInfo = cli.Command{
	Name:  "appinfo",
	Usage: "Show application template detail",
	Description: `
	This command obtains a detailed information about a specified application
	template. <app_name> parameter must contain the application template name.
	The <os_name> parameter must contain the name of the operating system template
	for which the template is designed. The <app_name> and <os_name> parameters
	together uniquely identify an application template. There could be multiple
	templates with the same name but designed for different operating systems.
`,
	Flags: CommonFlags,
	Action: func(c *cli.Context) {
		action(c, doApplicationInfo)
	},
}

var commandApplicationInstall = cli.Command{
	Name:  "appinstall",
	Usage: "Install application(s) into Container",
	Description: `
	This command installs application template (or multiple templates) into
	a Container. The application template must be compatible with the OS template
	installed in the target Container.

	You can specify multiple application templates in this command call.
`,
	Flags: CommonFlags,
	Action: func(c *cli.Context) {
		action(c, doApplicationInstall)
	},
}

var commandApplicationReset = cli.Command{
	Name:  "appreset",
	Usage: "Reset application(s) in Container",
	Description: `
	This command resets the installed application templates in a Container. It
	takes a list of application templates that you want installed in the Container.
	If the templates are already installed, they will be left untouched. If not,
	the will be installed. All other installed application templates will be
	removed from the server.
`,
	Flags: CommonFlags,
	Action: func(c *cli.Context) {
		action(c, doApplicationReset)
	},
}

var commandApplicationDelete = cli.Command{
	Name:  "appdelete",
	Usage: "Delete application from Container/Virtual machine",
	Description: `
	This command removes an application template from a Container. The <app_name>
	argument must contain the application name.
`,
	Flags: CommonFlags,
	Action: func(c *cli.Context) {
		action(c, doApplicationDelete)
	},
}

var commandOSList = cli.Command{
	Name:      "oslist",
	ShortName: "osls",
	Usage:     "List OS Templates used for creating a new server",
	Description: `
	This command obtains a list of the available operating system templates.

	An operating system template is a package which is used to create new server.
	It contains a particular operating system type and version (and software
	applications in some cases) together with necessary instructions and is used
	to preconfigure a server and install the operating system into it.

	When creating a server, use this command to obtain a list of the available OS
	templates, then choose the template of interest and use its name as an
	argument in the server creating command.

	The <os_name> argument is optional and may contain the OS template name. When
	it is included, only the information about the specified template will be
	retrieved.
`,
	Flags: append(CommonFlags, verboseFlag, noHeaderFlag),
	Action: func(c *cli.Context) {
		action(c, doOSList)
	},
}

func doApplicationList(c *cli.Context) {
	resp, err := client.SendRequest("GET", "/application-template", nil)
	assert(err)

	applist := lib.ApplicationList{}
	assert(xml.Unmarshal(resp.Body, &applist))
	if resp.StatusCode >= 400 {
		displayErrorAndExit(string(resp.Body))
	}

	outputResult(c, applist, func(format string) {
		tbl, err := prettytable.NewTable([]prettytable.Column{
			{Header: "ID", AlignRight: true},
			{Header: "NAME"},
			{Header: "FOROS"},
			{Header: "DESCRIPTION"},
		}...)
		assert(err)
		if c.Bool("no-header") {
			tbl.NoHeader = true
		}
		for _, e := range applist.ApplicationTemplate {
			tbl.AddRow(e.ID, e.Name, e.ForOS, e.Description)
		}
		tbl.Print()
	})
}

func doApplicationInfo(c *cli.Context) {
	if len(c.Args()) < 2 {
		displayWrongNumOfArgsAndExit(c)
	}
	appname := c.Args().Get(0)
	foros := c.Args().Get(1)

	resp, err := client.SendRequest("GET", "/application-template/"+appname+"/"+foros, nil)
	assert(err)
	if resp.StatusCode >= 400 {
		displayErrorAndExit(string(resp.Body))
	}

	app := lib.ApplicationTemplate{}
	assert(xml.Unmarshal(resp.Body, &app))

	outputResult(c, app, func(format string) {
		lib.PrintXMLStruct(app)
	})
}

func doApplicationInstall(c *cli.Context) {
	if len(c.Args()) < 2 {
		displayWrongNumOfArgsAndExit(c)
	}

	path := "/ve/" + c.Args().Get(0) + "/install"
	if len(c.Args()) == 2 {
		path += "/" + c.Args().Get(1)
	} else {
		for i, e := range c.Args()[1:] {
			if i == 0 {
				path += "?"
			} else {
				path += "&"
			}
			path += "name=" + e
		}
	}

	resp, err := client.SendRequest("PUT", path, nil)
	assert(err)

	switch resp.StatusCode {
	case 202:
		fmt.Println(string(resp.Body))
	default:
		displayErrorAndExit(string(resp.Body))
	}
}

func doApplicationReset(c *cli.Context) {
	if len(c.Args()) < 2 {
		displayWrongNumOfArgsAndExit(c)
	}

	path := "/ve/" + c.Args().Get(0) + "/application"
	for i, e := range c.Args()[1:] {
		if i == 0 {
			path += "?"
		} else {
			path += "&"
		}
		path += "name=" + e
	}

	resp, err := client.SendRequest("POST", path, nil)
	assert(err)

	switch resp.StatusCode {
	case 202:
		fmt.Println(string(resp.Body))
	default:
		displayErrorAndExit(string(resp.Body))
	}
}

func doApplicationDelete(c *cli.Context) {
	if len(c.Args()) < 2 {
		displayWrongNumOfArgsAndExit(c)
	}

	vename := c.Args().Get(0)
	appname := c.Args().Get(1)

	resp, err := client.SendRequest("DELETE", "/ve/"+vename+"/application/"+appname, nil)
	assert(err)

	switch resp.StatusCode {
	case 202:
		fmt.Println(string(resp.Body))
	default:
		displayErrorAndExit(string(resp.Body))
	}
}

func doOSList(c *cli.Context) {
	tmplName := ""
	if len(c.Args()) > 0 {
		tmplName = "/" + c.Args().Get(0)
	}

	resp, err := client.SendRequest("GET", "/template"+tmplName, nil)
	assert(err)
	if resp.StatusCode >= 400 {
		displayErrorAndExit(string(resp.Body))
	}

	if bytes.Contains(resp.Body, []byte("template-list")) {
		tmpls := lib.TemplateList{}
		assert(xml.Unmarshal(resp.Body, &tmpls))
		outputResult(c, tmpls, func(format string) {
			if c.Bool("verbose") {
				lib.PrintXMLStruct(tmpls)
			} else {
				tbl, err := prettytable.NewTable([]prettytable.Column{
					{Header: "TEMPLATE_NAME"}, {Header: "TECHNOLOGY"}, {Header: "TYPE"},
				}...)
				assert(err)
				if c.Bool("no-header") {
					tbl.NoHeader = true
				}
				for _, e := range tmpls.Template {
					tbl.AddRow(e.Name, e.Technology, e.OSType)
				}
				tbl.Print()
			}
		})
	} else {
		tmpl := lib.Template{}
		assert(xml.Unmarshal(resp.Body, &tmpl))
		outputResult(c, tmpl, func(format string) {
			if c.Bool("verbose") {
				lib.PrintXMLStruct(tmpl)
			} else {
				tbl, err := prettytable.NewTable([]prettytable.Column{
					{Header: "TEMPLATE_NAME"}, {Header: "TECHNOLOGY"}, {Header: "TYPE"},
				}...)
				assert(err)
				if c.Bool("no-header") {
					tbl.NoHeader = true
				}
				tbl.AddRow(tmpl.Name, tmpl.Technology, tmpl.OSType)
				tbl.Print()
			}
		})
	}
}
