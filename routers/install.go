package routers

import (
	"strings"

	"github.com/go-martini/martini"


	"github.com/coderoamer/shipper/modules/setting"
	"github.com/coderoamer/shipper/modules/log"
	"github.com/coderoamer/shipper/models"
)

func checkRunMode() {
	switch setting.Cfg.MustValue("", "RUN_MODE") {
	case "prod":
		martini.Env = martini.Prod
		setting.ProdMode = true
	case "test":
		martini.Env = martini.Test
	}
	log.Info("Run Mode: %s", strings.Title(martini.Env))
}

// GlobalInit is for global configuration reload-able.
func GlobalInit() {
	setting.NewConfigContext()
	log.Trace("Log path: %s", setting.LogRootPath)

	setting.NewServices()

	if setting.InstallLock {
		if err := models.NewEngine(); err != nil {
			log.Fatal("Fail to initialize ORM engine: %v", err)
		}

		models.HasEngine = true
	}

	if models.EnableSQLite3 {
		log.Info("SQLite3 Enabled")
	}
	checkRunMode()

}
