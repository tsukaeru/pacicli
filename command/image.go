package command

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/tatsushid/go-prettytable"
	"github.com/tsukaeru/pacicli/lib"
)

var commandImageList = cli.Command{
	Name:      "imglist",
	ShortName: "imgls",
	Usage:     "List images",
	Description: `
	This command obtains a list of the existing server images.
`,
	Flags: append(CommonFlags, noHeaderFlag),
	Action: func(c *cli.Context) {
		action(c, doImageList)
	},
}

var commandImageInfo = cli.Command{
	Name:  "imginfo",
	Usage: "Show image detail",
	Description: `
	This command obtains a detailed information for the specified server image.
`,
	Flags: CommonFlags,
	Action: func(c *cli.Context) {
		action(c, doImageInfo)
	},
}

var commandImageCreate = cli.Command{
	Name:  "imgcreate",
	Usage: "Create image from Container/Virtual machine",
	Description: `
	This command creates an image from an existing server. The server must be
	stopped before you attempt to create an image from it. The <server_name>
	argument must contain the source server name. The <image_name> must contain
	the desired image name.

	If you have multiple subscriptions, you have to specify the subscription ID
	by --subscription-id option. If not, it isn't required.
`,
	Flags: append(CommonFlags, subscriptionIDFlag),
	Action: func(c *cli.Context) {
		action(c, doImageCreate)
	},
}

var commandImageDelete = cli.Command{
	Name:  "imgdelete",
	Usage: "Delete image",
	Description: `
	This command deletes an existing server image.
`,
	Flags: CommonFlags,
	Action: func(c *cli.Context) {
		action(c, doImageDelete)
	},
}

func doImageList(c *cli.Context) {
	resp, err := client.SendRequest("GET", "/image", nil)
	assert(err)

	imglist := lib.ImageList{}
	assert(xml.Unmarshal(resp.Body, &imglist))

	outputResult(c, imglist, func(format string) {
		tbl, err := prettytable.NewTable([]prettytable.Column{
			{Header: "NAME"},
			{Header: "SIZE", AlignRight: true},
			{Header: "CREATED"},
			{Header: "SUBSCR_ID", AlignRight: true},
			{Header: "IMAGE_OF"},
			{Header: "DESCRIPTION"},
		}...)
		assert(err)
		if c.Bool("no-header") {
			tbl.NoHeader = true
		}
		for _, e := range imglist.ImageInfo {
			tbl.AddRow(e.Name, e.Size, e.Created, e.SubscriptionID, e.ImageOf, e.Description)
		}
		tbl.Print()
	})
}

func doImageInfo(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayWrongNumOfArgsAndExit(c)
	}
	imgname := c.Args().Get(0)

	resp, err := client.SendRequest("GET", "/image/"+imgname, nil)
	assert(err)
	if resp.StatusCode >= 400 {
		displayErrorAndExit(string(resp.Body))
	}

	img := lib.VeImage{}
	assert(xml.Unmarshal(resp.Body, &img))

	outputResult(c, img, func(format string) {
		lib.PrintXMLStruct(img)
	})
}

func doImageCreate(c *cli.Context) {
	if len(c.Args()) < 2 {
		displayWrongNumOfArgsAndExit(c)
	}
	vename := c.Args().Get(0)
	imgname := c.Args().Get(1)

	path := "/image/" + vename
	if c.Int("subscription-id") > 0 {
		path += "/" + strconv.Itoa(c.Int("subscription-id"))
	}
	path += "/create/" + imgname

	resp, err := client.SendRequest("POST", path, nil)
	assert(err)

	switch resp.StatusCode {
	case 202:
		fmt.Println(string(resp.Body))
	default:
		displayErrorAndExit(string(resp.Body))
	}
}

func doImageDelete(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayWrongNumOfArgsAndExit(c)
	}
	imgname := c.Args().Get(0)

	resp, err := client.SendRequest("DELETE", "/image/"+imgname, nil)
	assert(err)
	switch resp.StatusCode {
	case 202:
		fmt.Println(imgname, string(resp.Body))
	default:
		displayErrorAndExit(string(resp.Body))
	}
}
