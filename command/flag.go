package command

import (
	"github.com/codegangsta/cli"
	"github.com/tsukaeru/pacicli/lib"
)

var CommonFlags = []cli.Flag{
	configFileFlag, outputFlag,
}

var configFileFlag = cli.StringFlag{
	Name:   "config, c",
	Value:  "Pacifile",
	Usage:  "Specify config file path",
	EnvVar: "PACICLI_CONFIG",
}

var outputFlag = cli.StringFlag{
	Name:  "output, o",
	Value: "text",
	Usage: "Specify output format",
}

var verboseFlag = cli.BoolFlag{
	Name:  "verbose, v",
	Usage: "Verbose output",
}

var noHeaderFlag = cli.BoolFlag{
	Name:  "no-header, H",
	Usage: "Don't output column header",
}

var templateFlag = cli.StringFlag{
	Name:  "template, T",
	Usage: "Specify a template name",
}

var dropAppsFlag = cli.BoolFlag{
	Name:  "drop-apps, D",
	Usage: "Specify whether reinstalling applications into\n\tthe recreated server",
}

var descriptionFlag = cli.StringFlag{
	Name:  "description, desc",
	Usage: "Specify a server description",
}

var cpusFlag = cli.IntFlag{
	Name:  "cpus",
	Usage: "Specify a number of CPU cores",
}

var cpuPowerFlag = cli.IntFlag{
	Name:  "cpu-power",
	Usage: "Specify CPU clock rate in Mhz",
}

var ramSizeFlag = cli.IntFlag{
	Name:  "ram-size, ram",
	Usage: "Specify RAM size in MB",
}

var bandwidthFlag = cli.IntFlag{
	Name:  "bandwidth",
	Usage: "Specify bandwidth in kbps",
}

var addIPv4Flag = cli.IntFlag{
	Name:  "add-ipv4",
	Usage: "Specify a number of IPv4 addresses to be added.\n\tIt can't be used with --drop-ipv4",
}

var dropIPv4Flag = cli.StringSliceFlag{
	Name:  "drop-ipv4",
	Value: &cli.StringSlice{},
	Usage: "Specify IPv4 addresses to be removed.\n\tYou can use this option more than once\n\tand can't use with --add-ipv4",
}

var addIPv6Flag = cli.IntFlag{
	Name:  "add-ipv6",
	Usage: "Specify a number of IPv6 addresses to be added\n\tIt can't be used with --drop-ipv6",
}

var dropIPv6Flag = cli.StringSliceFlag{
	Name:  "drop-ipv6",
	Value: &cli.StringSlice{},
	Usage: "Specify IPv6 addresses to be removed.\n\tYou can use this option more than once\n\tand can't use with --add-ipv6",
}

var diskSizeFlag = cli.IntFlag{
	Name:  "disk-size",
	Usage: "Specify Disk size in GB",
}

var customNsFlag = cli.BoolFlag{
	Name:  "custom-ns",
	Usage: "Specify to keep modification to name server\n\tsetting in the server",
}

var noCustomNsFlag = cli.BoolFlag{
	Name:  "no-custom-ns",
	Usage: "Specify not to keep modification to name server\n\tsetting in the server",
}

var numRecordsFlag = cli.IntFlag{
	Name:  "num-records, n",
	Value: 10,
	Usage: "Specify a number of records API should return",
}

var fromDatetimeFlag = cli.StringFlag{
	Name:  "from, f",
	Usage: "Specify history start date and time in\n\t" + lib.ArgTimestampFormatStr() + " format",
}

var toDatetimeFlag = cli.StringFlag{
	Name:  "to, t",
	Usage: "Specify history end date and time in\n\t" + lib.ArgTimestampFormatStr() + " format",
}

var settingFlag = cli.StringFlag{
	Name:  "setting-file, s",
	Usage: "Specify a file path which contains server setting",
}

var subscriptionIDFlag = cli.IntFlag{
	Name:  "subscription-id, s",
	Usage: "Specify a subscription ID number",
}

var averagePeriodFlag = cli.IntFlag{
	Name:  "average-period",
	Usage: "Specify an interval in seconds for which to\n\tcalculate an average resource consumption value",
}

var tailFlag = cli.IntFlag{
	Name:  "tail",
	Usage: "Specify the time in seconds at the end of\n\tthe average-period for which the averageing\n\tshould NOT be performed",
}
