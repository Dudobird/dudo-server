package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/zhangmingkai4315/dudo-server/auth"
	"github.com/zhangmingkai4315/dudo-server/controllers"
)

func main() {
	router := mux.NewRouter()
	router.Use(auth.JWTAuthentication)
	router = controllers.Init(router)

	port := os.Getenv("server_port")
	if port == "" {
		port = "8080"
	}
	log.Printf("Info: server will listen at :%s\n", port)
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Panicln("Error:", err)
	}
}
