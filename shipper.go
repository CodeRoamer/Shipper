package main

import (
	"os"

	"github.com/codegangsta/cli"

	"github.com/coderoamer/shipper/cmd"
)

const APP_VER = "0.1.0"

func main() {
	app := cli.NewApp()
	app.Name = "Shipper"
	app.Usage = "Docker Web UI written in GO"
	app.Version = APP_VER
	app.Commands = []cli.Command{
		cmd.CmdWeb,
		cmd.CmdDump,
	}

	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Run(os.Args)
}
