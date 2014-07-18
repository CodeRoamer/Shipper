package cmd

import (
	"github.com/go-martini/martini"
	"github.com/codegangsta/cli"
	"net/http"
)

var CmdWeb = cli.Command{
	Name:  "web",
	Usage: "Start Shipper web server",
	Description: `Shipper web server is the only thing you need to run,
and it takes care of all the other things for you`,
	Action: runWeb,
	Flags:  []cli.Flag{},
}

func newMartini() *martini.ClassicMartini {
	r := martini.NewRouter()
	m := martini.New()
	m.Use(martini.Logger())
	m.Use(martini.Recovery())
	m.Use(martini.Static("public", martini.StaticOptions{}))
	m.MapTo(r, (*martini.Routes)(nil))
	m.Action(r.Handle)
	return &martini.ClassicMartini{m, r}
}


func runWeb(*cli.Context) {
	m := newMartini()

	m.Get("/", func() {
			println("hello world")
		})

	http.ListenAndServe("0.0.0.0:1212", m)
}
