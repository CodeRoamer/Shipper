package cmd

import (
	"fmt"
	"net/http"
	"html/template"

	"github.com/go-martini/martini"

	"github.com/martini-contrib/render"
	"github.com/martini-contrib/method"

	"github.com/codegangsta/cli"

	"github.com/coderoamer/shipper/routers"
	"github.com/coderoamer/shipper/modules/setting"
	"github.com/coderoamer/shipper/modules/log"
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
	return martini.Classic()
}


func runWeb(*cli.Context) {
	routers.GlobalInit()

	m := newMartini()

	m.Use(render.Renderer(render.Options{
		Directory: "templates", // Specify what path to load the templates from.
		Layout: "layout", // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"}, // Specify extensions to load for templates.
		Funcs: []template.FuncMap{}, // Specify helper function maps for templates to access.
		IndentJSON: true, // Output human readable JSON
		IndentXML: true, // Output human readable XML
	}))

	// checks for the X-HTTP-Method-Override header and uses it if the original request method is POST.
	m.Use(method.Override())

	m.Get("/", routers.Home)

	var err error
	listenAddr := fmt.Sprintf("%s:%s", setting.HttpAddr, setting.HttpPort)
	log.Info("Listen: %v://%s", setting.Protocol, listenAddr)

	switch setting.Protocol {
	case setting.HTTP:
		err = http.ListenAndServe(listenAddr, m)
	case setting.HTTPS:
		err = http.ListenAndServeTLS(listenAddr, setting.CertFile, setting.KeyFile, m)
	default:
		log.Fatal("Invalid protocol: %s", setting.Protocol)
	}

	if err != nil {
		log.Fatal("Fail to start server: %v", err)
	}
}
