package command

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/tatsushid/go-prettytable"
	"github.com/tsukaeru/pacicli/lib"
)

var commandLbList = cli.Command{
	Name:      "lblist",
	ShortName: "lbls",
	Usage:     "List load balancers",
	Description: `
	This command obtains a list of the available load balancers.
`,
	Flags: append(CommonFlags, noHeaderFlag),
	Action: func(c *cli.Context) {
		action(c, doLbList)
	},
}

var commandLbInfo = cli.Command{
	Name:  "lbinfo",
	Usage: "Show load balancer detail",
	Description: `
	This command obtains the information about a specified load balancer. The
	<lb_name> argument must contain the load balancer name.
`,
	Flags: append(CommonFlags, verboseFlag),
	Action: func(c *cli.Context) {
		action(c, doLbInfo)
	},
}

var commandLbHistory = cli.Command{
	Name:      "lbhistory",
	ShortName: "lbhist",
	Usage:     "Show load balancer history",
	Description: `
	The command obtains a load balancer history. The <lb_name> argument must contain
	the name of the load balancer for which to retrieve the history and -n option
	must be used to specify the number of records to be included in the result set.
`,
	Flags: append(CommonFlags, numRecordsFlag, verboseFlag, noHeaderFlag),
	Action: func(c *cli.Context) {
		action(c, doLbHistory)
	},
}

var commandLbCreate = cli.Command{
	Name:  "lbcreate",
	Usage: "Create load balancer",
	Description: `
	This command creates a load balancer. The <lb_name> argument must contain
	the user defined load balancer name.

	If you have multiple subscriptions, you have to specify the subscription ID
	by --subscription-id option. If not, it isn't required.
`,
	Flags: append(CommonFlags, subscriptionIDFlag),
	Action: func(c *cli.Context) {
		action(c, doLbCreate)
	},
}

var commandLbRestart = cli.Command{
	Name:  "lbrestart",
	Usage: "Restart load balancer",
	Description: `
	This command restarts a load balancer.
`,
	Flags: CommonFlags,
	Action: func(c *cli.Context) {
		action(c, doLbRestart)
	},
}

var commandLbDelete = cli.Command{
	Name:  "lbdelete",
	Usage: "Delete load balancer",
	Description: `
	This command deletes an existing load balancer. If there are servers attached
	to a load balancer, the deletion fails. In such a case, the servers have to be
	detached first.
`,
	Flags: CommonFlags,
	Action: func(c *cli.Context) {
		action(c, doLbDelete)
	},
}

var commandLbAttach = cli.Command{
	Name:  "lbattach",
	Usage: "Attach Container/Virtual machine to load balancer",
	Description: `
	The command attaches a server to a load balancer. Once this request is
	completed, the server load will be managed by the specified load balancer.
	The <lb_name> and <server_name> must contain the load balancer and the server
	names respectively.
`,
	Flags: CommonFlags,
	Action: func(c *cli.Context) {
		action(c, doLbAttachDetach)
	},
}

var commandLbDetach = cli.Command{
	Name:  "lbdetach",
	Usage: "Detach Container/Virtual machine to load balancer",
	Description: `
	The command detaches a server from a load balancer. The <lb_name> and
	<server_name> must contain the load balancer and the server names respectively.
`,
	Flags: CommonFlags,
	Action: func(c *cli.Context) {
		action(c, doLbAttachDetach)
	},
}

func doLbList(c *cli.Context) {
	resp, err := client.SendRequest("GET", "/load-balancer", nil)
	assert(err)

	lblist := lib.LbList{}
	assert(xml.Unmarshal(resp.Body, &lblist))

	outputResult(c, lblist, func(format string) {
		tbl, err := prettytable.NewTable([]prettytable.Column{
			{Header: "NAME"}, {Header: "STATE"}, {Header: "SUBSCR_ID", AlignRight: true},
		}...)
		assert(err)
		if c.Bool("no-header") {
			tbl.NoHeader = true
		}
		for _, e := range lblist.LoadBalancer {
			tbl.AddRow(e.Name, e.State, e.SubscriptionID)
		}
		tbl.Print()
	})
}

func doLbInfo(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayWrongNumOfArgsAndExit(c)
	}
	lbname := c.Args().Get(0)

	resp, err := client.SendRequest("GET", "/load-balancer/"+lbname, nil)
	assert(err)
	if resp.StatusCode >= 400 {
		displayErrorAndExit(string(resp.Body))
	}

	lb := lib.LoadBalancer{}
	assert(xml.Unmarshal(resp.Body, &lb))

	outputResult(c, lb, func(format string) {
		if c.Bool("verbose") {
			lib.PrintXMLStruct(lb)
		} else {
			var publicIP lib.IPAddr
			if len(lb.Network.PublicIP) > 0 {
				publicIP = lb.Network.PublicIP[0].Address
			}
			fmt.Println("LOAD BALANCER INFO")
			fmt.Printf("             Name: %s\n", lb.Name)
			fmt.Printf("  Subscription ID: %d\n", lb.SubscriptionID)
			fmt.Printf("Public IP address: %s\n", publicIP)
			fmt.Printf("           Status: %s\n\n", lb.State)
			fmt.Println("BALANCED SERVERS")

			tbl, err := prettytable.NewTable([]prettytable.Column{
				{Header: "NAME"},
				{Header: "IPADDR"},
			}...)
			assert(err)
			for _, e := range lb.UsedBy {
				tbl.AddRow(e.VeName, e.IP)
			}
			tbl.Print()
		}
	})
}

func doLbHistory(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayWrongNumOfArgsAndExit(c)
	}
	lbname := c.Args().Get(0)

	if c.Int("num-records") <= 0 {
		displayErrorAndExit("This command must be used with --num-records flag and its argument. Please see '" + c.App.Name + " help " + c.Command.Name + "'")
	}
	numRecord := strconv.Itoa(c.Int("num-records"))

	resp, err := client.SendRequest("GET", "/load-balancer/"+lbname+"/history/"+numRecord, nil)
	assert(err)
	if resp.StatusCode >= 400 {
		displayErrorAndExit(string(resp.Body))
	}

	hst := lib.VeHistory{}
	assert(xml.Unmarshal(resp.Body, &hst))

	outputResult(c, hst, func(format string) {
		if c.Bool("verbose") {
			lib.PrintXMLStruct(hst)
		} else {
			tbl, err := prettytable.NewTable([]prettytable.Column{
				{Header: "DATETIME"},
				{Header: "CPU", AlignRight: true},
				{Header: "MEMORY", AlignRight: true},
				{Header: "DISK", AlignRight: true},
				{Header: "BANDWIDTH", AlignRight: true},
				{Header: "PUB_IPS", AlignRight: true},
				{Header: "STATUS"},
			}...)
			assert(err)
			if c.Bool("no-header") {
				tbl.NoHeader = true
			}
			for _, e := range hst.VeSnapshot {
				ts, _ := e.EventTimestamp.MarshalText()
				tbl.AddRow(ts, e.CPU, e.RAM, e.LocalDisk, e.Bandwidth, e.NoOfPublicIP, e.State)
			}
			tbl.Print()
		}
	})
}

func doLbCreate(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayWrongNumOfArgsAndExit(c)
	}
	lbname := c.Args().Get(0)

	path := "/load-balancer"
	if c.Int("subscription-id") > 0 {
		path += "/" + strconv.Itoa(c.Int("subscription-id"))
	}
	path += "/create/" + lbname

	resp, err := client.SendRequest("POST", path, nil)
	assert(err)

	if resp.StatusCode >= 400 {
		displayErrorAndExit(string(resp.Body))
	}

	pwd := lib.PasswordResponse{}
	assert(xml.Unmarshal(resp.Body, &pwd))

	outputResult(c, pwd, func(format string) {
		lib.PrintXMLStruct(pwd)
	})
}

func doLbRestart(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayWrongNumOfArgsAndExit(c)
	}
	lbname := c.Args().Get(0)

	resp, err := client.SendRequest("PUT", "/load-balancer/"+lbname+"/restart", nil)
	assert(err)
	switch resp.StatusCode {
	case 202:
		fmt.Println(lbname, string(resp.Body))
	default:
		displayErrorAndExit(string(resp.Body))
	}
}

func doLbDelete(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayWrongNumOfArgsAndExit(c)
	}
	lbname := c.Args().Get(0)

	resp, err := client.SendRequest("DELETE", "/load-balancer/"+lbname, nil)
	assert(err)
	switch resp.StatusCode {
	case 202:
		fmt.Println(lbname, string(resp.Body))
	default:
		displayErrorAndExit(string(resp.Body))
	}
}

func doLbAttachDetach(c *cli.Context) {
	if len(c.Args()) < 2 {
		displayWrongNumOfArgsAndExit(c)
	}
	lbname := c.Args().Get(0)
	vename := c.Args().Get(1)

	method := "POST"
	if c.Command.Name == "lbdetach" {
		method = "DELETE"
	}
	resp, err := client.SendRequest(method, "/load-balancer/"+lbname+"/"+vename, nil)
	assert(err)

	switch resp.StatusCode {
	case 202:
		fmt.Println(string(resp.Body))
	default:
		displayErrorAndExit(string(resp.Body))
	}
}
