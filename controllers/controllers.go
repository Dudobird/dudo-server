package controllers

import (
	"github.com/gorilla/mux"
)

// Init will registe all controllers to router and return it
// we will call this method in main function
func Init(router *mux.Router) *mux.Router {
	router.HandleFunc("/api/user/new", CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", Login).Methods("POST")

	return router
}
