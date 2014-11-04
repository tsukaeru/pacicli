package command

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/tatsushid/go-prettytable"
	"github.com/tsukaeru/pacicli/lib"
)

var commandAutoscale = cli.Command{
	Name:  "autoscale",
	Usage: "Show auto scaling rules of Container/Virtual machine",
	Description: `
	This command obtains auto scaling rules for the specified server
`,
	Flags: CommonFlags,
	Action: func(c *cli.Context) {
		action(c, doAutoscale)
	},
}

var commandAutoscaleCreate = cli.Command{
	Name:  "autoscale-create",
	Usage: "Create auto scaling rules for Container/Virtual machine",
	Description: `
	This command creates auto scaling rules for the specified server using rules
	specified in Pacicli or --setting-file flag argument file
`,
	Flags: append(CommonFlags, settingFlag),
	Action: func(c *cli.Context) {
		action(c, doAutoscaleCreateUpdate)
	},
}

var commandAutoscaleUpdate = cli.Command{
	Name:  "autoscale-update",
	Usage: "Update auto scaling rules of Container/Virtual machine",
	Description: `
	This command updates existing auto scaling rules for the specified server
	using rules specified in Pacicli or --setting-file flag argument file
`,
	Flags: append(CommonFlags, settingFlag),
	Action: func(c *cli.Context) {
		action(c, doAutoscaleCreateUpdate)
	},
}

var commandAutoscaleDrop = cli.Command{
	Name:  "autoscale-drop",
	Usage: "Drop auto scaling rules from Container/Virtual machine",
	Description: `
	This command drops auto scaling rules from the specified server
`,
	Flags: CommonFlags,
	Action: func(c *cli.Context) {
		action(c, doAutoscaleDrop)
	},
}

var commandAutoscaleHistory = cli.Command{
	Name:  "autoscale-history",
	Usage: "Get auto scaling history of Container/Virtual machine",
	Description: `
	This command obtains auto scaling history. It must be used with a pair of
	--from and --to flags datetime arguments or --num-records flag argument.

	With --from and --to pair, --average-period and --tail flags can be used.

	--average-period is used to specify an interval (in seconds) for which to
	calculate an average resource consumption value. For example, if the argument
	is included and contains a value of 900 seconds, only an average value of every
	consecutive 15-minute period will be included in the response. If the argument
	is not included, the entire set of the actual values will be returned. Please
	note that the network traffic values are calculated as maximum over the
	specified period. All other statistics are calculated as average values.

	--tail is used to specify the time (in seconds) at the end of the
	average-period for which the averaging should NOT be performed. For example,
	if --average-period 3600 and --tail 900, one average value (or the maximum
	value for network traffic) will be included in the response for the first 45
	minutes of every consecutive one-hour period. For the last 15 minutes of every
	hour, the response will contain a complete set of the actual values.

	You should also be aware that if multiple command call are executed against
	the same server and are expecting to receive the complete statistical data,
	the CPU, memory, and network load on the server side may increase drastically,
	which may significantly slow down the processing of the command call. Using
	the averaging approach, you can avoid this potential problem.
`,
	Flags: append(CommonFlags, numRecordsFlag, fromDatetimeFlag, toDatetimeFlag, averagePeriodFlag, tailFlag, verboseFlag),
	Action: func(c *cli.Context) {
		action(c, doAutoscaleHistory)
	},
}

func doAutoscale(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayWrongNumOfArgsAndExit(c)
	}
	vename := c.Args().Get(0)

	resp, err := client.SendRequest("GET", "/ve/"+vename+"/autoscale", nil)
	assert(err)
	if resp.StatusCode >= 400 {
		displayErrorAndExit(string(resp.Body))
	}

	autoscale := lib.Autoscale{}
	assert(xml.Unmarshal(resp.Body, &autoscale))

	outputResult(c, autoscale, func(format string) {
		lib.PrintXMLStruct(autoscale)
	})
}

func doAutoscaleCreateUpdate(c *cli.Context) {
	method := "PUT"
	if c.Command.Name == "autoscale-create" {
		method = "POST"
	}

	if len(c.Args()) == 0 {
		displayWrongNumOfArgsAndExit(c)
	}
	vename := c.Args().Get(0)

	var data lib.AutoscaleData
	if len(c.String("setting-file")) > 0 {
		assert(lib.LoadConfig(c.String("setting-file"), &data))
	} else {
		if s, ok := conf.Servers[vename]; ok && len(s.AutoscaleRule) > 0 {
			data = lib.AutoscaleData{AutoscaleRule: s.AutoscaleRule}
		} else {
			displayErrorAndExit("Couldn't find Autoscale rules for '" + vename + "'")
		}
	}
	var b bytes.Buffer
	assert(xml.NewEncoder(&b).Encode(data))

	resp, err := client.SendRequest(method, "/ve/"+vename+"/autoscale", &b)
	assert(err)

	if resp.StatusCode != 200 {
		displayErrorAndExit(string(resp.Body))
	}

	autoscale := lib.Autoscale{}
	assert(xml.Unmarshal(resp.Body, &autoscale))

	outputResult(c, autoscale, func(format string) {
		lib.PrintXMLStruct(autoscale)
	})
}

func doAutoscaleDrop(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayWrongNumOfArgsAndExit(c)
	}
	vename := c.Args().Get(0)

	resp, err := client.SendRequest("DELETE", "/ve/"+vename+"/autoscale", nil)
	assert(err)

	switch resp.StatusCode {
	case 200:
		fmt.Println(string(resp.Body))
	default:
		displayErrorAndExit(string(resp.Body))
	}
}

func doAutoscaleHistory(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayWrongNumOfArgsAndExit(c)
	}
	vename := c.Args().Get(0)

	path := "/ve/" + vename + "/autoscale/history/"
	if len(c.String("from")) > 0 && len(c.String("to")) > 0 {
		_, err := lib.ParseArgTimestampFormat(c.String("from"))
		assert(err, "'from' arg value must be in "+lib.ArgTimestampFormatStr()+" format")

		_, err = lib.ParseArgTimestampFormat(c.String("to"))
		assert(err, "'to' arg value must be in "+lib.ArgTimestampFormatStr()+" format")

		path += c.String("from") + "/" + c.String("to")
		var q []string
		if c.Int("average-period") > 0 {
			q = append(q, "average-period="+strconv.Itoa(c.Int("average-period")))
		}
		if c.Int("tail") > 0 {
			q = append(q, "tail="+strconv.Itoa(c.Int("tail")))
		}
		if len(q) > 0 {
			path += "?" + strings.Join(q, "&")
		}
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

	hst := lib.ResourceConsumptionAndAutoscaleHistory{}
	assert(xml.Unmarshal(resp.Body, &hst))
	assert(err)

	outputResult(c, hst, func(format string) {
		if c.Bool("verbose") {
			lib.PrintXMLStruct(hst)
		} else {
			if len(hst.AutoscaleRule) > 0 {
				fmt.Println("AUTOSCALE RULE HISTORY")
				tbl, err := prettytable.NewTable([]prettytable.Column{
					{Header: "METRIC"},
					{Header: "VERSION", AlignRight: true},
					{Header: "UPDATED"},
					{Header: "DELIVERED"},
					{Header: "DELIVERED-OK"},
					{Header: "MIGRATION"},
					{Header: "RESTART"},
					{Header: "MIN", AlignRight: true},
					{Header: "MAX", AlignRight: true},
					{Header: "STEP", AlignRight: true},
					{Header: "UP_THRES", AlignRight: true},
					{Header: "UP_PERIOD", AlignRight: true},
					{Header: "DOWN_THRES", AlignRight: true},
					{Header: "DOWN_PERIOD", AlignRight: true},
				}...)
				assert(err)
				if c.Bool("no-header") {
					tbl.NoHeader = true
				}
				for _, e := range hst.AutoscaleRule {
					tbl.AddRow(
						e.Metric,
						*e.Version,
						*e.Updated,
						*e.UpdateDelivered,
						*e.UpdateDeliveredOk,
						*e.AllowMigration,
						*e.AllowRestart,
						e.Limits.Min,
						e.Limits.Max,
						e.Limits.Step,
						*e.Thresholds.Up.Threshold,
						e.Thresholds.Up.Period,
						*e.Thresholds.Down.Threshold,
						e.Thresholds.Down.Period,
					)
				}
				tbl.Print()
				fmt.Println()
			}

			fmt.Println("RESOURCE CONSUMPTION")
			tbl, err := prettytable.NewTable([]prettytable.Column{
				{Header: "CPU_USAGE", AlignRight: true},
				{Header: "RAM_USAGE", AlignRight: true},
				{Header: "PRIV_IN", AlignRight: true},
				{Header: "PRIV_OUT", AlignRight: true},
				{Header: "PUB_IN", AlignRight: true},
				{Header: "PUB_OUT", AlignRight: true},
				{Header: "DATETIME"},
				{Header: "CPU", AlignRight: true},
				{Header: "RAM", AlignRight: true},
				{Header: "BANDWIDTH", AlignRight: true},
			}...)
			assert(err)
			if c.Bool("no-header") {
				tbl.NoHeader = true
			}
			for _, e := range hst.ResourceConsumptionSample {
				tbl.AddRow(
					e.CPUUsage,
					e.RAMUsage,
					e.PrivateIncomingTraffic,
					e.PrivateOutgoingTraffic,
					e.PublicIncomingTraffic,
					e.PublicOutgoingTraffic,
					e.PaciTimestamp,
					e.CPU,
					e.RAM,
					e.Bandwidth,
				)
			}
			tbl.Print()
		}
	})
}
