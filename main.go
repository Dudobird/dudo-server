package main

import (
	"flag"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/zhangmingkai4315/dudo-server/config"
	"github.com/zhangmingkai4315/dudo-server/routers"
)

var configFile string

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	flag.StringVar(&configFile, "c", "config.toml", "config file path")
}

func main() {
	flag.Parse()
	config := config.LoadConfig(configFile)
	router := routers.LoadRouters()
	hostAndPort := config.Application.ListenAt
	log.Println("server will listen at ", hostAndPort)
	err := http.ListenAndServe(hostAndPort, router)
	if err != nil {
		log.Fatal(err)
	}
}
