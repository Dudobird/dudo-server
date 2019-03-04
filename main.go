package main

import (
	"flag"
	"os"

	"github.com/Dudobird/dudo-server/core"
	"github.com/Dudobird/dudo-server/routers"
	log "github.com/sirupsen/logrus"
)

var configFile string

func init() {
	log.SetReportCaller(true)
	log.SetLevel(log.InfoLevel)
	log.SetOutput(os.Stdout)
	flag.StringVar(&configFile, "c", "config.toml", "config file path")
}

func main() {
	flag.Parse()
	app := core.NewApp(configFile)
	router, err := routers.LoadRouters()
	if err != nil {
		return
	}
	app.Router = router
	app.Run()
}
