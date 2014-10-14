package command

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/codegangsta/cli"
	"github.com/tatsushid/go-prettytable"
	"github.com/tsukaeru/pacicli/lib"
)

type Command struct {
	cli.Command
	Synopsis string
}

var Commands = []cli.Command{
	commandList,
	commandStart,
	commandStop,
	commandCreate,
	commandCreateFromImage,
	commandClone,
	commandRecreate,
	commandModify,
	commandResetPassword,
	commandInfo,
	commandHistory,
	commandUsage,
	commandDelete,
	commandFirewallList,
	commandFirewallCreate,
	commandFirewallModify,
	commandFirewallDelete,
	commandBackupScheduleSet,
	commandBackupScheduleRemove,
	commandBackup,
	commandBackupList,
	commandBackupRestore,
	commandBackupInfo,
	commandBackupDelete,
	commandAutoscale,
	commandAutoscaleCreate,
	commandAutoscaleUpdate,
	commandAutoscaleDrop,
	commandAutoscaleHistory,
	commandApplicationList,
	//commandApplicationInfo,
	commandApplicationInstall,
	commandApplicationReset,
	commandApplicationDelete,
	commandImageList,
	commandImageInfo,
	commandImageCreate,
	commandImageDelete,
	commandLbList,
	commandLbInfo,
	commandLbHistory,
	commandLbCreate,
	commandLbRestart,
	commandLbDelete,
	commandLbAttach,
	commandLbDetach,
	commandOSList,
	commandBackupSchedule,
	commandInitiatingVnc,
}

var commandSynopsisses = map[string]string{
	"list":                   "[options]",
	"start":                  "<server_name> [options]",
	"stop":                   "<server_name> [options]",
	"create":                 "<server_name> [options]",
	"create-from-image":      "<server_name> <image_name> [options]",
	"clone":                  "<src_server_name> <dst_server_name> [options]",
	"recreate":               "<server_name> [options]",
	"modify":                 "<server_name> [options]",
	"reset-password":         "<server_name> [options]",
	"info":                   "<server_name> [options]",
	"history":                "<server_name> {-f <from> -t <to> | -n <num>} [options]",
	"usage":                  "<server_name> -f <from> -t <to> [options]",
	"delete":                 "<server_name> [options]",
	"vnc":                    "<server_name> [options]",
	"fwlist":                 "<server_name> [options]",
	"fwcreate":               "<server_name> [options]",
	"fwmodify":               "<server_name> [options]",
	"fwdelete":               "<server_name> [options]",
	"backup-schedule-set":    "<server_name> <schedule_name> [options]",
	"backup-schedule-remove": "<server_name> [options]",
	"backup":                 "<server_name> [options]",
	"backup-list":            "<server_name> -f <from> -t <to> [options]",
	"backup-restore":         "<server_name> <backup_id> [options]",
	"backup-info":            "<server_name> <backup_id> [options]",
	"backup-delete":          "<server_name> <backup_id> [options]",
	"autoscale":              "<server_name> [options]",
	"autoscale-create":       "<server_name> [options]",
	"autoscale-update":       "<server_name> [options]",
	"autoscale-drop":         "<server_name> [options]",
	"autoscale-history":      "<server_name> {-f <from> -t <to> | -n <num>} [options]",
	"applist":                "[options]",
	"appinfo":                "<app_name> <os_name> [options]",
	"appinstall":             "<server_name> <app_name ...> [options]",
	"appreset":               "<server_name> <app_name ...> [options]",
	"appdelete":              "<server_name> <app_name> [options]",
	"imglist":                "[options]",
	"imginfo":                "<image_name> [options]",
	"imgcreate":              "<server_name> <image_name> [options]",
	"imgdelete":              "<image_name> [options]",
	"lblist":                 "[options]",
	"lbinfo":                 "<lb_name> [options]",
	"lbhistory":              "<lb_name> -n <num> [options]",
	"lbcreate":               "<lb_name> [options]",
	"lbrestart":              "<lb_name> [options]",
	"lbdelete":               "<lb_name> [options]",
	"lbattach":               "<lb_name> <server_name> [options]",
	"lbdetach":               "<lb_name> <server_name> [options]",
	"oslist":                 "[<os_name>] [options]",
	"backup-schedule":        "[options]",
}

const (
	jsonPrefix      = ""
	jsonIndent      = "  "
	columnSeparator = "   "
)

func assert(err error, v ...interface{}) {
	if err != nil {
		a := []interface{}{err}
		if len(v) > 0 {
			a = append(a, ": ")
			a = append(a, v...)
		}
		displayErrorAndExit(a...)
	}
}

func displayErrorAndExit(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}

func displayWrongNumOfArgsAndExit(c *cli.Context) {
	cmdStr := c.App.Name + " help " + c.Command.Name
	displayErrorAndExit("The command was used with the wrong number of arguments. Please see '" + cmdStr + "' result")
}

func getBackupID(id string) string {
	if !strings.HasPrefix(id, "{") {
		id = "{" + id
	}
	if !strings.HasSuffix(id, "}") {
		id += "}"
	}
	return id
}

func outputResult(c *cli.Context, v interface{}, defaultFn func(format string)) error {
	f := strings.ToLower(c.String("output"))
	switch f {
	case "json":
		b, err := json.MarshalIndent(v, jsonPrefix, jsonIndent)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	case "toml":
		return toml.NewEncoder(os.Stdout).Encode(v)
	default:
		defaultFn(f)
	}
	return nil
}

var (
	conf   lib.Config
	client *lib.Client
)

func action(c *cli.Context, fn func(c *cli.Context)) {
	if len(c.String("config")) > 0 {
		err := lib.LoadConfig(c.String("config"), &conf)
		assert(err)
		if len(conf.BaseURL) == 0 || len(conf.Username) == 0 || len(conf.Password) == 0 {
			displayErrorAndExit("Invalid config data. BaseURL, Username and Password must be correctly specified in a config file")
		}
		client = lib.NewClient(conf.BaseURL, conf.Username, conf.Password)
		prettytable.Separator = columnSeparator
		fn(c)
	} else {
		displayErrorAndExit("Config path is empty. It must be specified to use this command.\nPlease see '" + c.App.Name + " help' result")
	}
}
