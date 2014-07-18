package cmd

import (
	"log"

	"github.com/codegangsta/cli"

	"github.com/coderoamer/shipper/models"
	"github.com/coderoamer/shipper/modules/setting"
)

var CmdDump = cli.Command{
	Name:  "dump",
	Usage: "Dump Shipper files and database",
	Description: `Dump compresses all related files and database into zip file.
It can be used for backup and capture Shipper server image to send to maintainer`,
	Action: runDump,
	Flags:  []cli.Flag{},
}

func runDump(*cli.Context) {
	setting.NewConfigContext()
	models.LoadModelsConfig()
	models.SetEngine()

	log.Println("Finish dumping!")
}
