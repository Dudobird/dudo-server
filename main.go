package main

import (
	"log"
	"net/http"
	"os"

	"github.com/zhangmingkai4315/dudo-server/routers"
)

func main() {
	router := routers.InitRouters()
	port := os.Getenv("server_port")
	if port == "" {
		port = "8080"
	}
	log.Printf("Info: server will listen at :%s\n", port)
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Panicln("Error:", err)
	}
}
