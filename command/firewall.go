package command

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/tatsushid/go-prettytable"
	"github.com/tsukaeru/pacicli/lib"
)

var commandFirewallList = cli.Command{
	Name:      "fwlist",
	ShortName: "fwls",
	Usage:     "List firewall rules",
	Description: `
	The command obtains a list of existing firewall rules for the specified server.
	The <server_name> must contain the server name.
`,
	Flags: append(CommonFlags, noHeaderFlag),
	Action: func(c *cli.Context) {
		action(c, doFirewallList)
	},
}

var commandFirewallCreate = cli.Command{
	Name:  "fwcreate",
	Usage: "Create firewall rules for Container/Virtual machine",
	Description: `
	This command creates firewall rules. The <server_name> must contain the server
	name and firewall rules must be defined in Pacicli or --setting-file flag
	argument file.

	You have to use 'fwmodify' command instead of this when firewall rules has
	already existed.
`,
	Flags: append(CommonFlags, settingFlag),
	Action: func(c *cli.Context) {
		action(c, doFirewallCreateModify)
	},
}

var commandFirewallModify = cli.Command{
	Name:  "fwmodify",
	Usage: "Modify firewall rules of Container/Virtual machine",
	Description: `
	This command modifies existing firewall rules. It replaces all existing rules
	with the new ones. To keep existing rules and add more, first obtain the list
	of the existing rules, then add new rules to it and use the complete list as
	an input.

	The <server_name> must contain the server name and firewall rules must be
	defined in Pacicli or --setting-file flag argument file.
`,
	Flags: append(CommonFlags, settingFlag),
	Action: func(c *cli.Context) {
		action(c, doFirewallCreateModify)
	},
}

var commandFirewallDelete = cli.Command{
	Name:  "fwdelete",
	Usage: "Delete firewall rules from Container/Virtual machine",
	Description: `
	The command deletes all existing firewall rules. To delete a specific rule,
	retrieve all existing rules, modify the result set as needed and then use it
	as an argument of the 'fwmodify' command.
`,
	Flags: CommonFlags,
	Action: func(c *cli.Context) {
		action(c, doFirewallDelete)
	},
}

func doFirewallList(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayWrongNumOfArgsAndExit(c)
	}
	vename := c.Args().Get(0)

	resp, err := client.SendRequest("GET", "/ve/"+vename+"/firewall", nil)
	assert(err)
	if resp.StatusCode >= 400 {
		displayErrorAndExit(string(resp.Body))
	}

	fwlist := lib.Firewall{}
	assert(xml.Unmarshal(resp.Body, &fwlist))

	outputResult(c, fwlist, func(format string) {
		tbl, err := prettytable.NewTable([]prettytable.Column{
			{Header: "ID", AlignRight: true},
			{Header: "NAME"},
			{Header: "PROTOCOL"},
			{Header: "LOCAL_PORT", AlignRight: true},
			{Header: "REMOTE_PORT", AlignRight: true},
			{Header: "REMOTE_NET"},
		}...)
		assert(err)
		if c.Bool("no-header") {
			tbl.NoHeader = true
		}
		for _, e := range fwlist.Rule {
			var ra lib.IPAddr
			if len(e.RemoteNet) > 0 {
				ra = e.RemoteNet[0]
			}
			tbl.AddRow(e.ID, e.Name, e.Protocol, e.LocalPort, e.RemotePort, ra)
			if len(e.RemoteNet) > 1 {
				for _, a := range e.RemoteNet[1:] {
					tbl.AddRow("", "", "", "", "", a)
				}
			}
		}
		tbl.Print()
	})
}

func doFirewallCreateModify(c *cli.Context) {
	method := "PUT"
	if c.Command.Name == "fwcreate" {
		method = "POST"
	}

	if len(c.Args()) == 0 {
		displayWrongNumOfArgsAndExit(c)
	}
	vename := c.Args().Get(0)

	var fw lib.Firewall
	if len(c.String("setting-file")) > 0 {
		assert(lib.LoadConfig(c.String("setting-file"), &fw))
	} else {
		if s, ok := conf.Servers[vename]; ok && len(s.Firewall.Rule) > 0 {
			fw = s.Firewall
		} else {
			displayErrorAndExit("Couldn't find Firewall rules for '" + vename + "'")
		}
	}

	var b bytes.Buffer
	assert(xml.NewEncoder(&b).Encode(fw))

	resp, err := client.SendRequest(method, "/ve/"+vename+"/firewall", &b)
	assert(err)

	switch resp.StatusCode {
	case 200:
		fmt.Println(string(resp.Body))
	default:
		displayErrorAndExit(string(resp.Body))
	}
}

func doFirewallDelete(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayWrongNumOfArgsAndExit(c)
	}
	vename := c.Args().Get(0)

	resp, err := client.SendRequest("DELETE", "/ve/"+vename+"/firewall", nil)
	assert(err)

	switch resp.StatusCode {
	case 200:
		fmt.Println(string(resp.Body))
	default:
		displayErrorAndExit(string(resp.Body))
	}
}
