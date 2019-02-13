package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/zhangmingkai4315/dudo-server/config"
	"github.com/zhangmingkai4315/dudo-server/routers"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "c", "config.toml", "config file path")
}

func main() {
	flag.Parse()
	config := config.LoadConfig(configFile)
	router := routers.InitRouters()
	hostAndPort := config.Application.ListenAt
	log.Println("Info: server will listen at ", hostAndPort)
	err := http.ListenAndServe(hostAndPort, router)
	if err != nil {
		log.Panicln("Error:", err)
	}
}
