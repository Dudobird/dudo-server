package routers

import (
	"fmt"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/zhangmingkai4315/dudo-server/auth"
	"github.com/zhangmingkai4315/dudo-server/controllers"
)

// LoadRouters will registe all controllers to router and return it
// we will call this method in main function
func LoadRouters() (router *mux.Router, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("panic: %v", r)
			}
		}
	}()
	router = mux.NewRouter()
	router.Use(auth.JWTAuthentication)
	router.HandleFunc("/api/auth/signup", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/auth/signin", controllers.Login).Methods("POST")

	// router.HandleFunc("/api/user/{id}", controllers.CreateAccount).Methods("PUT")
	// router.HandleFunc("/api/user/{id}", controllers.CreateAccount).Methods("DELETE")
	// router.HandleFunc("/api/user/{id}", controllers.CreateAccount).Methods("GET")
	router.HandleFunc("/api/users", controllers.CreateAccount).Methods("GET")

	log.Infoln("load api routers success")
	// // static files
	// router.Handle("/", http.FileServer(http.Dir("../frontend/")))
	// router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("../frontend/static/"))))
	return
}
