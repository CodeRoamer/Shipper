package main

import (
	"os"
	"runtime"

	"github.com/codegangsta/cli"

	"github.com/coderoamer/shipper/cmd"
	"github.com/coderoamer/shipper/modules/setting"
)

const APP_VER = "0.1.0"

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	setting.AppVer = APP_VER
}

func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true
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
