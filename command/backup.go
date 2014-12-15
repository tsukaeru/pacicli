package command

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/tatsushid/go-prettytable"
	"github.com/tsukaeru/pacicli/lib"
)

var commandBackupScheduleSet = cli.Command{
	Name:  "backup-schedule-set",
	Usage: "Assign a backup schedule to a Container/Virtual machine",
	Description: `
	This command assigns a backup schedule to the specified server. Backup
	schedules are created and configured by the system administrator and define
	when and how often the server backups will be performed. Backup schedules also
	define the maximum number of incremental backups and a maximum number of
	backups to keep on the backup server.

	The <schedule_name> argument must contain the name of a predefined backup
	schedule. You can obtain the list of the existing schedules using the
	'backup-schedule' command. You cannot create your own backup schedules.
`,
	Flags: CommonFlags,
	Action: func(c *cli.Context) {
		action(c, doBackupScheduleSet)
	},
}

var commandBackupScheduleRemove = cli.Command{
	Name:  "backup-schedule-remove",
	Usage: "Remove a backup schedule from a Container/Virtual machine",
	Description: `
	This command removes a backup schedule from the specified server.
`,
	Flags: CommonFlags,
	Action: func(c *cli.Context) {
		action(c, doBackupScheduleRemove)
	},
}

var commandBackup = cli.Command{
	Name:  "backup",
	Usage: "Perform an on demand backup of Container/Virtual machine",
	Description: `
	This command performs an on-demand backup of the specified server.
`,
	Flags: CommonFlags,
	Action: func(c *cli.Context) {
		action(c, doBackup)
	},
}

var commandBackupList = cli.Command{
	Name:      "backup-list",
	ShortName: "bkpls",
	Usage:     "List backups of a Container/Virtual machine",
	Description: `
	This command lists available backups for the specified server. The --from and
	--to flags arguments must be used with it to specify datetime interval for
	which to retrieve the backups.
`,
	Flags: append(CommonFlags, fromDatetimeFlag, toDatetimeFlag, verboseFlag),
	Action: func(c *cli.Context) {
		action(c, doBackupList)
	},
}

var commandBackupRestore = cli.Command{
	Name:  "backup-restore",
	Usage: "Restore a Container/Virtual machine from a backup",
	Description: `
	The command restores a specified server from a specified backup. A server
	must be stopped in order to perform a restore operation. The <backup_id>
	argument must contain the backup ID. Please see 'backup-list' command for the
	information on how to obtain the list of the available backups and their IDs.

	Please note that the complete backup ID string must be specified as the
	command argument, including curly brackets and any other leading and trailing
	characters (if any).
`,
	Flags: CommonFlags,
	Action: func(c *cli.Context) {
		action(c, doBackupRestore)
	},
}

var commandBackupInfo = cli.Command{
	Name:  "backup-info",
	Usage: "Show backup detail",
	Description: `
	This command obtains the information about the specified backup. The
	<backup_id> must contain a valid backup ID. (Please see 'backup-list' command)
`,
	Flags: CommonFlags,
	Action: func(c *cli.Context) {
		action(c, doBackupInfo)
	},
}

var commandBackupDelete = cli.Command{
	Name:  "backup-delete",
	Usage: "Delete a backup of a Container/Virtual machine",
	Description: `
	This command deletes a specified backup. Please note that you can only delete
	an on-demand backup. Scheduled backups cannot be deleted by the user.

	The <backup_id> argument must contain a valid backup ID (Please see
	'backup-list' command)
`,
	Flags: CommonFlags,
	Action: func(c *cli.Context) {
		action(c, doBackupDelete)
	},
}

var commandBackupSchedule = cli.Command{
	Name:  "backup-schedule",
	Usage: "List backup schedules",
	Description: `
	This command obtains the list of the available backup schedules. Backup
	schedules are created by system administrator. If you would like to perform
	server backups on a regular basis, you can obtain the list of the available
	schedules using this command, then choose a schedule that suits your needs and
	specify its name when configuring your server.
`,
	Flags: append(CommonFlags, noHeaderFlag),
	Action: func(c *cli.Context) {
		action(c, doBackupSchedule)
	},
}

func doBackupScheduleSet(c *cli.Context) {
	if len(c.Args()) < 2 {
		displayWrongNumOfArgsAndExit(c)
	}
	vename := c.Args().Get(0)
	schedule := c.Args().Get(1)

	resp, err := client.SendRequest("PUT", "/ve/"+vename+"/schedule/"+schedule, nil)
	assert(err)

	switch resp.StatusCode {
	case 202:
		fmt.Println(string(resp.Body))
	default:
		displayErrorAndExit(string(resp.Body))
	}
}

func doBackupScheduleRemove(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayWrongNumOfArgsAndExit(c)
	}
	vename := c.Args().Get(0)

	resp, err := client.SendRequest("PUT", "/ve/"+vename+"/nobackup/", nil)
	assert(err)

	switch resp.StatusCode {
	case 202:
		fmt.Println(string(resp.Body))
	default:
		displayErrorAndExit(string(resp.Body))
	}
}

func doBackup(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayWrongNumOfArgsAndExit(c)
	}
	vename := c.Args().Get(0)

	resp, err := client.SendRequest("POST", "/ve/"+vename+"/backup", nil)
	assert(err)

	switch resp.StatusCode {
	case 202:
		fmt.Println(string(resp.Body))
	default:
		displayErrorAndExit(string(resp.Body))
	}
}

func doBackupList(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayWrongNumOfArgsAndExit(c)
	}
	vename := c.Args().Get(0)

	if len(c.String("from")) == 0 || len(c.String("to")) == 0 {
		displayErrorAndExit("This command must be used with a pair of --from and --to flags arguments. Please see '" + c.App.Name + " help " + c.Command.Name + "'")
	}
	path := "/ve/" + vename + "/backups/"

	from, err := lib.ParseArgTimestampFormat(c.String("from"))
	assert(err, "'from' arg value must be in "+lib.ArgTimestampFormatStr()+" format")

	to, err := lib.ParseArgTimestampFormat(c.String("to"))
	assert(err, "'to' arg value must be in "+lib.ArgTimestampFormatStr()+" format")

	path += c.String("from") + "/" + c.String("to")

	resp, err := client.SendRequest("GET", path, nil)
	assert(err)
	if resp.StatusCode >= 400 {
		displayErrorAndExit(string(resp.Body))
	}

	backups := lib.VeBackups{}
	assert(xml.Unmarshal(resp.Body, &backups))

	outputResult(c, backups, func(format string) {
		if c.Bool("verbose") {
			lib.PrintXMLStruct(backups)
		} else {
			fmt.Println("BACKUP LIST")
			fmt.Printf("Server: %s\n", vename)
			fmt.Printf("  From: %s\n", from.Format(lib.DataTimestampFormat))
			fmt.Printf("    To: %s\n\n", to.Format(lib.DataTimestampFormat))

			tbl, err := prettytable.NewTable([]prettytable.Column{
				{Header: "ID"},
				{Header: "SCHEDULE"},
				{Header: "START"},
				{Header: "END"},
				{Header: "RESULT", AlignRight: true},
				{Header: "SIZE(GB)", AlignRight: true},
				{Header: "NODE"},
				{Header: "DESCRIPTION"},
			}...)
			assert(err)
			for _, e := range backups.Backup {
				schedule := "-"
				if len(e.ScheduleName) > 0 {
					schedule = e.ScheduleName
				}
				result := "fail"
				if e.Successful == true {
					result = "ok"
				}
				size := strconv.FormatFloat(float64(e.BackupSize)/(1<<30), 'f', 3, 64)
				tbl.AddRow(e.CloudBackupID, schedule, e.Started, e.Ended, result, size, e.BackupNodeName, e.Description)
			}
			tbl.Print()
		}
	})
}

func doBackupRestore(c *cli.Context) {
	if len(c.Args()) < 2 {
		displayWrongNumOfArgsAndExit(c)
	}
	vename := c.Args().Get(0)
	backupid := getBackupID(c.Args().Get(1))

	resp, err := client.SendRequest("PUT", "/ve/"+vename+"/restore/"+backupid, nil)
	assert(err)

	switch resp.StatusCode {
	case 202:
		fmt.Println(string(resp.Body))
	default:
		displayErrorAndExit(string(resp.Body))
	}
}

func doBackupInfo(c *cli.Context) {
	if len(c.Args()) < 2 {
		displayWrongNumOfArgsAndExit(c)
	}
	vename := c.Args().Get(0)
	backupid := getBackupID(c.Args().Get(1))

	resp, err := client.SendRequest("GET", "/ve/"+vename+"/backup/"+backupid, nil)
	assert(err)
	if resp.StatusCode >= 400 {
		displayErrorAndExit(string(resp.Body))
	}

	backup := lib.Backup{}
	assert(xml.Unmarshal(resp.Body, &backup))

	outputResult(c, backup, func(format string) {
		lib.PrintXMLStruct(backup)
	})
}

func doBackupDelete(c *cli.Context) {
	if len(c.Args()) < 2 {
		displayWrongNumOfArgsAndExit(c)
	}
	vename := c.Args().Get(0)
	backupid := getBackupID(c.Args().Get(1))

	resp, err := client.SendRequest("DELETE", "/ve/"+vename+"/backup/"+backupid, nil)
	assert(err)

	switch resp.StatusCode {
	case 202:
		fmt.Println(string(resp.Body))
	default:
		displayErrorAndExit(string(resp.Body))
	}
}

func doBackupSchedule(c *cli.Context) {
	resp, err := client.SendRequest("GET", "/schedule", nil)
	assert(err)
	if resp.StatusCode >= 400 {
		displayErrorAndExit(string(resp.Body))
	}

	backups := lib.BackupScheduleList{}
	assert(xml.Unmarshal(resp.Body, &backups))

	outputResult(c, backups, func(format string) {
		tbl, err := prettytable.NewTable([]prettytable.Column{
			{Header: "ID", AlignRight: true},
			{Header: "NAME"},
			{Header: "DESCRIPTION"},
			{Header: "ENABLED", AlignRight: true},
			{Header: "KEEP", AlignRight: true},
			{Header: "INCREMENTAL", AlignRight: true},
		}...)
		assert(err)
		if c.Bool("no-header") {
			tbl.NoHeader = true
		}
		for _, e := range backups.BackupSchedule {
			tbl.AddRow(e.ID, e.Name, e.Description, e.Enabled, e.BackupsToKeep, e.NoOfIncremental)
		}
		tbl.Print()
	})
}
