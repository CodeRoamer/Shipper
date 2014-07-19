package cmd

import (
	"os/exec"
	"bytes"

	"github.com/codegangsta/cli"

	"github.com/coderoamer/shipper/modules/log"
)

var CmdTest = cli.Command{
	Name:  "test",
	Usage: "Test Against All the Packages in Shipper",
	Description: `The command will run the tests in all shipper's modules, use -b or --benchmark to specify the extra need for running benchmark test`,
	Action: runTest,
	Flags:  []cli.Flag{
		cli.StringFlag{Name: "benchmark, b", Value: "false", Usage: "run benchmark test or not"},
	},
}

func runTest(ctx *cli.Context) {
	log.Trace("Running Test...")

	var bench bool = false

	if ctx.String("benchmark") == "true" {
		bench = true
	}

	conditionalTest("github.com/coderoamer/shipper", bench)
	conditionalTest("github.com/coderoamer/shipper/models", bench)
	conditionalTest("github.com/coderoamer/shipper/modules/base", bench)
	conditionalTest("github.com/coderoamer/shipper/modules/middleware", bench)

	log.Trace("Test Complete!")
}

func conditionalTest(pkg string, bench bool) {
	var cmd *exec.Cmd
	if bench {
		cmd = exec.Command("go","test","-bench=.",pkg)
	} else {
		cmd = exec.Command("go","test",pkg)
	}

	var (
		stdOut bytes.Buffer
		stdErr bytes.Buffer
	)

	cmd.Stderr = &stdErr
	cmd.Stdout = &stdOut

	err := cmd.Run()
	if err != nil {
		log.Warn("Test against %s has errors:", pkg)
		if len(stdOut.String())!=0 {
			log.Error("%s", stdOut.String())
		}
		if len(stdErr.String())!=0 {
			log.Error("%s", stdErr.String())
		}
	} else {
		log.Info("%s", stdOut.String())
	}

}
