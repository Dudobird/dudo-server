package routers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zhangmingkai4315/dudo-server/auth"
	"github.com/zhangmingkai4315/dudo-server/controllers"
)

// InitRouters will registe all controllers to router and return it
// we will call this method in main function
func InitRouters() *mux.Router {
	router := mux.NewRouter()
	router.Use(auth.JWTAuthentication)
	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Login).Methods("POST")

	// static files
	router.Handle("/", http.FileServer(http.Dir("../frontend/")))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("../frontend/static/"))))
	return router
}
