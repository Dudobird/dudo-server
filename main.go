package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/Dudobird/dudo-server/core"
)

var configFile string

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	flag.StringVar(&configFile, "c", "config.toml", "config file path")
}

func main() {
	flag.Parse()
	app := core.NewApp(configFile)
	app.Run()
}
