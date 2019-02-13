package routers

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/zhangmingkai4315/dudo-server/auth"
	"github.com/zhangmingkai4315/dudo-server/controllers"
)

// LoadRouters will registe all controllers to router and return it
// we will call this method in main function
func LoadRouters() *mux.Router {
	router := mux.NewRouter()
	router.Use(auth.JWTAuthentication)
	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Login).Methods("POST")
	log.Infoln("Load routers success")
	// // static files
	// router.Handle("/", http.FileServer(http.Dir("../frontend/")))
	// router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("../frontend/static/"))))
	return router
}
