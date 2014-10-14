package command

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"github.com/tatsushid/go-prettytable"
	"github.com/tsukaeru/pacicli/lib"
)

var commandList = cli.Command{
	Name:      "list",
	ShortName: "ls",
	Usage:     "List Containers/Virtual machines",
	Description: `
	The command obtains the list of servers owned by the current user. You can
	use --subscription-id option to list only the servers that belong to a specific
	subscription.
`,
	Flags: append(CommonFlags, subscriptionIDFlag, noHeaderFlag),
	Action: func(c *cli.Context) {
		action(c, doList)
	},
}

var commandStart = cli.Command{
	Name:  "start",
	Usage: "Start Container/Virtual machine",
	Description: `
	This command starts a specific server.
`,
	Flags: CommonFlags,
	Action: func(c *cli.Context) {
		action(c, doStartStop)
	},
}

var commandStop = cli.Command{
	Name:  "stop",
	Usage: "Stop Container/Virtual machine",
	Description: `
	This command stops a specific server.
`,
	Flags: CommonFlags,
	Action: func(c *cli.Context) {
		action(c, doStartStop)
	},
}

var commandCreate = cli.Command{
	Name:  "create",
	Usage: "Create Container/Virtual machine",
	Description: `
	This command creates a new server. To create it, a server --setting-file option
	must be used and all server properties are correctly defined in it.

	The administrator password for the new server will be automatically generated
	and displayed as a command result.

	If you have multiple subscriptions, you have to specify the subscription ID
	in the setting file. If not, it isn't required.
`,
	Flags: append(CommonFlags, settingFlag),
	Action: func(c *cli.Context) {
		action(c, doCreate)
	},
}

var commandCreateFromImage = cli.Command{
	Name:  "create-from-image",
	Usage: "Create Container/Virtual machine from image",
	Description: `
	This command creates a server from an existing image. The <server_name>
	argument must contain the new server name (user-defined). The <image_name>
	argument must contain the source image name. Please see also 'imglist' command
	and other related command's helps.

	The administrator password for the new server will be automatically generated
	and displayed as a command result.

	If you have multiple subscriptions, you have to specify the subscription ID
	by --subscription-id option. If not, it isn't required.
`,
	Flags: append(CommonFlags, subscriptionIDFlag),
	Action: func(c *cli.Context) {
		action(c, doCreateFromImage)
	},
}

var commandClone = cli.Command{
	Name:  "clone",
	Usage: "Clone Container/Virtual machine",
	Description: `
	This command creates a clone of an existing server. The <src_server_name>
	argument must contain the name of the source server. The <dst_server_name>
	argument must contain the new server name (user-defined). The IP addresses,
	gateway and DNS setting for the new server will be set automatically. The
	rest of the server configuration will be inherited from the source server.

	The administrator password for the new server will be automatically generated
	and displayed as a command result.

	If you have multiple subscriptions, you have to specify the subscription ID
	by --subscription-id option. If not, it isn't required.
`,
	Flags: append(CommonFlags, subscriptionIDFlag),
	Action: func(c *cli.Context) {
		action(c, doClone)
	},
}

var commandRecreate = cli.Command{
	Name:  "recreate",
	Usage: "Recreate Container/Virtual machine",
	Description: `
	The command recreates an existing server with the possibility of using a
	different OS template. This essentially deletes an existing server and then
	creates a new server keeping the original server specifications. Instead of
	deleting a server and creating a new one using separate operations, you can
	do the same using this single command.

	The new server will have the same resources (CPU, RAM, bandwidth, disk size
	and storage type) as the previous one. The new server will also have the same
	IP addresses, which cannot be guaranteed when deleting and then creating in
	two separate steps. The new server will also have the same auto scale and
	firewall rules and backup schedule if they were configured for the original
	server. If the original server was attached to a load balancer, the new one
	will be attached as well.

	The optional --template flag argument is used to specify an OS template. If the
	argument is not specified, the server will be recreated using the same OS
	template as the original. To obtain the list of the available templates, use
	the 'oslist' command

	The optional --drop-apps flag argument is used to specify whether to reinstall
	the applications in the new server. If the argument is specified, the
	applications installed in the original server will also be dropped from the new
	server. If you don't specify the argument, all application that are installed in
	the original server will be installed in the new one.
`,
	Flags: append(CommonFlags, templateFlag, dropAppsFlag),
	Action: func(c *cli.Context) {
		action(c, doRecreate)
	},
}

var commandModify = cli.Command{
	Name:  "modify",
	Usage: "Modify Container/Virtual machine configuration",
	Description: `
	This command modifies the configuration of an existing server. Please note that
	this command cannot be used to modify the server backup schedule. Please see
	'backup-schedule-set' command and related.

	To modify it, --setting-file or other configuration options must be used and
	please remenber you can't use --add-ipv(4|6) and --drop-ipv(4|6) at the same
	time.
`,
	Flags: append(CommonFlags, settingFlag, cpusFlag, cpuPowerFlag, ramSizeFlag, bandwidthFlag, addIPv4Flag, dropIPv4Flag, addIPv6Flag, dropIPv6Flag, diskSizeFlag, customNsFlag, noCustomNsFlag),
	Action: func(c *cli.Context) {
		action(c, doModify)
	},
}

var commandResetPassword = cli.Command{
	Name:  "reset-passwd",
	Usage: "Reset Container/Virtual machine password",
	Description: `
	This command resets the server administrator password. The new password will
	be automatically generated.
`,
	Flags: CommonFlags,
	Action: func(c *cli.Context) {
		action(c, doResetPassword)
	},
}

var commandInfo = cli.Command{
	Name:  "info",
	Usage: "Show Container/Virtual machine detail",
	Description: `
	This command obtains the information about the specified server. The
	<server_name> argument must contain the server name.
`,
	Flags: CommonFlags,
	Action: func(c *cli.Context) {
		action(c, doInfo)
	},
}

var commandHistory = cli.Command{
	Name:  "history",
	Usage: "Show Container/Virtual machine history",
	Description: `
	This command obtains the modification history for the specified server. Every
	time a server configuration is changed, a snapshot of the configuration is
	taken and saved. This command allows to retrieve these records and use them
	for statistics.
	
	This command must be used with a pair of --from and --to flags datetime
	arguments or --num-records flag argument.
`,
	Flags: append(CommonFlags, numRecordsFlag, fromDatetimeFlag, toDatetimeFlag, verboseFlag, noHeaderFlag),
	Action: func(c *cli.Context) {
		action(c, doHistory)
	},
}

var commandUsage = cli.Command{
	Name:  "usage",
	Usage: "Show Container/Virtual machine resource usage report",
	Description: `
	This command obtains the usage information for the specified server. To specify
	the datetime interval, it must be used with a pair of --from and --to flags
	datetime arguments.
`,
	Flags: append(CommonFlags, fromDatetimeFlag, toDatetimeFlag, verboseFlag),
	Action: func(c *cli.Context) {
		action(c, doUsage)
	},
}

var commandDelete = cli.Command{
	Name:  "delete",
	Usage: "Delete Container/Virtual machine",
	Description: `
	This command permanently deletes a server. Please note that you can only delete
	a fully stopped server. If a server is in a transition state (stopping,
	starting, a disk is being attached to it, etc.), it cannot be deleted.
`,
	Flags: CommonFlags,
	Action: func(c *cli.Context) {
		action(c, doDelete)
	},
}

var commandInitiatingVnc = cli.Command{
	Name:  "vnc",
	Usage: "Initiating VNC Session to Container/Virtual machine",
	Description: `
	This command initializes VNC console setting of the specified server. It returns
	an auto generated password to connect to the server via VNC client software. It
	is also needed that the IP address and the port information to connect. You can
	obtain them using 'info' command after the command. These will be included in
	the 'console' part of the 'info' command result.
`,
	Flags: CommonFlags,
	Action: func(c *cli.Context) {
		action(c, doInitiatingVnc)
	},
}

func doList(c *cli.Context) {
	path := "/ve"
	if c.Int("subscription-id") > 0 {
		path += "?subscription=" + strconv.Itoa(c.Int("subscription-id"))
	}
	resp, err := client.SendRequest("GET", path, nil)
	assert(err)

	velist := lib.VeList{}
	assert(xml.Unmarshal(resp.Body, &velist))
	if resp.StatusCode >= 400 {
		displayErrorAndExit(string(resp.Body))
	}

	tbl, err := prettytable.NewTable([]prettytable.Column{
		{Header: "ID", AlignRight: true},
		{Header: "NAME"},
		{Header: "HOSTNAME"},
		{Header: "STATE"},
		{Header: "SUBSCR_ID", AlignRight: true},
	}...)
	assert(err)

	outputResult(c, velist, func(format string) {
		if c.Bool("no-header") {
			tbl.NoHeader = true
		}
		for _, e := range velist.VeInfo {
			tbl.AddRow(e.ID, e.Name, e.Hostname, e.State, e.SubscriptionID)
		}
		tbl.Print()
	})
}

func doStartStop(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayWrongNumOfArgsAndExit(c)
	}
	vename := c.Args().Get(0)

	resp, err := client.SendRequest("PUT", "/ve/"+vename+"/"+c.Command.Name, nil)
	assert(err)
	switch resp.StatusCode {
	case 202:
		fmt.Println(vename, string(resp.Body))
	case 304:
		s := "started"
		if c.Command.Name == "stop" {
			s = "stopped"
		}
		displayErrorAndExit(vename, "has already", s)
	default:
		displayErrorAndExit(string(resp.Body))
	}
}

func doCreate(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayWrongNumOfArgsAndExit(c)
	}
	vename := c.Args().Get(0)

	var ve lib.CreateVe
	if len(c.String("setting-file")) > 0 {
		assert(lib.LoadConfig(c.String("setting-file"), &ve))
		ve.Name = vename
	} else {
		if s, ok := conf.Servers[vename]; ok && s.Spec != nil {
			ve = *s.Spec
		} else {
			cli.ShowCommandHelp(c, c.Command.Name)
			os.Exit(1)
		}
	}

	if len(ve.Hostname) == 0 {
		ve.Hostname = ve.Name
	}

	var b bytes.Buffer
	assert(xml.NewEncoder(&b).Encode(ve))

	resp, err := client.SendRequest("POST", "/ve/", &b)
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

func doCreateFromImage(c *cli.Context) {
	if len(c.Args()) < 2 {
		displayWrongNumOfArgsAndExit(c)
	}
	vename := c.Args().Get(0)
	imgname := c.Args().Get(1)

	path := "/ve"
	if c.Int("subscription-id") > 0 {
		path += "/" + strconv.Itoa(c.Int("subscription-id"))
	}
	path += "/" + vename + "/from/" + imgname

	resp, err := client.SendRequest("POST", path, nil)
	assert(err)
	switch resp.StatusCode {
	case 202:
		fmt.Println(vename, string(resp.Body))
	default:
		displayErrorAndExit(string(resp.Body))
	}
}

func doClone(c *cli.Context) {
	if len(c.Args()) < 2 {
		displayWrongNumOfArgsAndExit(c)
	}
	srcve := c.Args().Get(0)
	destve := c.Args().Get(1)

	path := "/ve/" + srcve + "/clone-to/" + destve
	if c.Int("subscription-id") > 0 {
		path += "/for/" + strconv.Itoa(c.Int("subscription-id"))
	}

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

func doRecreate(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayWrongNumOfArgsAndExit(c)
	}
	vename := c.Args().Get(0)

	path := "/ve/" + vename + "/recreate"
	var q []string
	if len(c.String("template")) > 0 {
		q = append(q, "template="+c.String("template"))
	}
	if c.Bool("drop-apps") {
		q = append(q, "drop-apps=true")
	}
	if len(q) > 0 {
		path += "?" + strings.Join(q, "&")
	}

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

func doModify(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayWrongNumOfArgsAndExit(c)
	}
	vename := c.Args().Get(0)

	var ve lib.ReconfigureVe
	if len(c.String("setting-file")) > 0 {
		assert(lib.LoadConfig(c.String("setting-file"), &ve))
	}
	if len(c.String("description")) > 0 {
		ve.Description = c.String("description")
	}
	if c.Int("cpus") > 0 {
		if ve.ChangeCPU == nil {
			ve.ChangeCPU = new(lib.ChangeCPU)
		}
		ve.ChangeCPU.Number = c.Int("cpus")
	}
	if c.Int("cpu-power") > 0 {
		if ve.ChangeCPU == nil {
			ve.ChangeCPU = new(lib.ChangeCPU)
		}
		ve.ChangeCPU.Power = c.Int("cpu-power")
	}
	if c.Int("ram-size") > 0 {
		ve.RAMSize = c.Int("ram-size")
	}
	if c.Int("bandwidth") > 0 {
		ve.Bandwidth = c.Int("bandwidth")
	}
	if c.Int("add-ipv4") > 0 {
		if ve.ReconfigureIPv4 == nil {
			ve.ReconfigureIPv4 = new(lib.ReconfigureIP)
		}
		if ve.ReconfigureIPv4.AddIP == nil {
			ve.ReconfigureIPv4.AddIP = new(lib.AddIP)
		}
		ve.ReconfigureIPv4.AddIP.Number = c.Int("add-ipv4")
	}
	if len(c.StringSlice("drop-ipv4")) > 0 {
		if ve.ReconfigureIPv4 == nil {
			ve.ReconfigureIPv4 = new(lib.ReconfigureIP)
		}
		if ve.ReconfigureIPv4.DropIP == nil {
			ve.ReconfigureIPv4.DropIP = new(lib.DropIP)
		}
		for _, addr := range c.StringSlice("drop-ipv4") {
			a, err := lib.NewIPAddr(addr)
			assert(err)
			ve.ReconfigureIPv4.DropIP.IP = append(ve.ReconfigureIPv4.DropIP.IP, *a)
		}
	}
	if c.Int("add-ipv6") > 0 {
		if ve.ReconfigureIPv4 == nil {
			ve.ReconfigureIPv6 = new(lib.ReconfigureIP)
		}
		if ve.ReconfigureIPv6.AddIP == nil {
			ve.ReconfigureIPv6.AddIP = new(lib.AddIP)
		}
		ve.ReconfigureIPv6.AddIP.Number = c.Int("add-ipv6")
	}
	if len(c.StringSlice("drop-ipv6")) > 0 {
		if ve.ReconfigureIPv6 == nil {
			ve.ReconfigureIPv6 = new(lib.ReconfigureIP)
		}
		if ve.ReconfigureIPv6.DropIP == nil {
			ve.ReconfigureIPv6.DropIP = new(lib.DropIP)
		}
		for _, addr := range c.StringSlice("drop-ipv6") {
			a, err := lib.NewIPAddr(addr)
			assert(err)
			ve.ReconfigureIPv6.DropIP.IP = append(ve.ReconfigureIPv6.DropIP.IP, *a)
		}
	}
	if c.Int("disk-size") > 0 {
		ve.PrimaryDiskSize = c.Int("disk-size")
	}
	if c.Bool("custom-ns") {
		if ve.CustomNs == nil {
			ve.CustomNs = new(int)
		}
		*(ve.CustomNs) = 1
	}
	if c.Bool("no-custom-ns") {
		if ve.CustomNs == nil {
			ve.CustomNs = new(int)
		}
		*(ve.CustomNs) = 0
	}

	if ve == (lib.ReconfigureVe{}) {
		displayErrorAndExit("There is no modification parameter. Please see '" + c.App.Name + " help " + c.Command.Name)
	}

	if ve.ReconfigureIPv4 != nil && ve.ReconfigureIPv4.AddIP != nil && ve.ReconfigureIPv4.DropIP != nil {
		displayErrorAndExit("Invalid modification setting. Both ReconfigureIPv4.AddIP and DropIP can't be specified at same time")
	}
	if ve.ReconfigureIPv6 != nil && ve.ReconfigureIPv6.AddIP != nil && ve.ReconfigureIPv6.DropIP != nil {
		displayErrorAndExit("Invalid modification setting. Both ReconfigureIPv6.AddIP and DropIP can't be specified at same time")
	}

	var b bytes.Buffer
	assert(xml.NewEncoder(&b).Encode(ve))

	resp, err := client.SendRequest("PUT", "/ve/"+vename, &b)
	assert(err)
	switch resp.StatusCode {
	case 202:
		fmt.Println(vename, string(resp.Body))
	default:
		displayErrorAndExit(string(resp.Body))
	}
}

func doResetPassword(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayWrongNumOfArgsAndExit(c)
	}
	vename := c.Args().Get(0)

	resp, err := client.SendRequest("POST", "/ve/"+vename+"/reset-password", nil)
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

func doInfo(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayWrongNumOfArgsAndExit(c)
	}
	vename := c.Args().Get(0)

	resp, err := client.SendRequest("GET", "/ve/"+vename, nil)
	assert(err)
	if resp.StatusCode >= 400 {
		displayErrorAndExit(string(resp.Body))
	}

	ve := lib.Ve{}
	assert(xml.Unmarshal(resp.Body, &ve))

	outputResult(c, ve, func(format string) {
		lib.PrintXMLStruct(ve)
	})
}

func doHistory(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayWrongNumOfArgsAndExit(c)
	}
	vename := c.Args().Get(0)

	path := "/ve/" + vename + "/history/"
	if len(c.String("from")) > 0 && len(c.String("to")) > 0 {
		_, err := time.Parse(lib.ArgTimestampFormat, c.String("from"))
		assert(err, "'from' arg value must be in a '"+lib.ArgTimestampFormat+"' format")

		_, err = time.Parse(lib.ArgTimestampFormat, c.String("to"))
		assert(err, "'to' arg value must be in a '"+lib.ArgTimestampFormat+"' format")

		path += c.String("from") + "/" + c.String("to")
	} else if c.Int("num-records") > 0 {
		path += strconv.Itoa(c.Int("num-records"))
	} else {
		displayErrorAndExit("This command must be used with a pair of --from and --to flags arguments or --num-records flag argument. Please see '" + c.App.Name + " help " + c.Command.Name + "'")
	}

	resp, err := client.SendRequest("GET", path, nil)
	assert(err)
	if resp.StatusCode >= 400 {
		displayErrorAndExit(string(resp.Body))
	}

	hst := lib.VeHistory{}
	assert(xml.Unmarshal(resp.Body, &hst))
	assert(err)

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

func doUsage(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayWrongNumOfArgsAndExit(c)
	}
	vename := c.Args().Get(0)

	if len(c.String("from")) == 0 || len(c.String("to")) == 0 {
		displayErrorAndExit("This command must be used with a pair of --from and --to flags arguments. Please see '" + c.App.Name + " help " + c.Command.Name + "'")
	}
	path := "/ve/" + vename + "/usage/"

	from, err := time.Parse(lib.ArgTimestampFormat, c.String("from"))
	assert(err, "'from' arg value must be in a '"+lib.ArgTimestampFormat+"' format")

	to, err := time.Parse(lib.ArgTimestampFormat, c.String("to"))
	assert(err, "'to' arg value must be in a '"+lib.ArgTimestampFormat+"' format")

	path += c.String("from") + "/" + c.String("to")

	resp, err := client.SendRequest("GET", path, nil)
	assert(err)
	if resp.StatusCode >= 400 {
		displayErrorAndExit(string(resp.Body))
	}

	usage := lib.VeResourceUsageReport{}
	assert(xml.Unmarshal(resp.Body, &usage))

	outputResult(c, usage, func(format string) {
		if c.Bool("verbose") {
			lib.PrintXMLStruct(usage)
		} else {
			fmt.Println("RESOURCE USAGE REPORT")
			fmt.Printf("Server: %s\n", usage.VeName)
			fmt.Printf("  From: %s\n", from.Format(lib.DataTimestampFormat))
			fmt.Printf("    To: %s\n\n", to.Format(lib.DataTimestampFormat))

			tbl, err := prettytable.NewTable([]prettytable.Column{
				{Header: "RESOURCE_TYPE"}, {Header: "USAGE", AlignRight: true},
			}...)
			assert(err)
			for _, e := range usage.ResourceUsage {
				name := e.ResourceType
				if len(e.ResourceUsageType) > 0 {
					name += "(" + e.ResourceUsageType + ")"
				}
				tbl.AddRow(name, e.Value)
			}
			for _, e := range usage.VeTraffic {
				tbl.AddRow(e.TrafficType, e.Used)
			}
			tbl.Print()
		}
	})
}

func doDelete(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayWrongNumOfArgsAndExit(c)
	}
	vename := c.Args().Get(0)

	resp, err := client.SendRequest("DELETE", "/ve/"+vename, nil)
	assert(err)
	switch resp.StatusCode {
	case 202:
		fmt.Println(vename, string(resp.Body))
	default:
		displayErrorAndExit(string(resp.Body))
	}
}

func doInitiatingVnc(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayWrongNumOfArgsAndExit(c)
	}
	vename := c.Args().Get(0)

	resp, err := client.SendRequest("POST", "/ve/"+vename+"/console", nil)
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
