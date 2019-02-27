package routers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Dudobird/dudo-server/auth"
	"github.com/Dudobird/dudo-server/controllers"
	"github.com/Dudobird/dudo-server/core"
	"github.com/Dudobird/dudo-server/utils"
	log "github.com/sirupsen/logrus"
)

const appContextKey = utils.ContextToken("App")

func appBindMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app := core.GetApp()
		ctx := context.WithValue(r.Context(), appContextKey, app)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

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
	router.Use(appBindMiddleware)
	router.HandleFunc("/api/auth/signup", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/api/auth/signin", controllers.Login).Methods("POST")
	router.HandleFunc("/api/auth/logout", controllers.Logout).Methods("GET")
	router.HandleFunc("/api/auth/password", controllers.UpdatePassword).Methods("UPDATE")

	router.HandleFunc("/api/storages", controllers.CreateFiles).Methods("POST")
	router.HandleFunc("/api/storages", controllers.GetTopLevelFiles).Methods("GET")
	router.HandleFunc("/api/storage/{id}", controllers.GetCurrentFile).Methods("GET")
	router.HandleFunc("/api/storage/{id}", controllers.DeleteFiles).Methods("DELETE")
	router.HandleFunc("/api/storage/{id}/subfiles", controllers.ListSubFiles).Methods("GET")

	router.HandleFunc("/api/upload/{parentID}", controllers.UploadFiles).Methods("POST")
	router.HandleFunc("/api/download/{id}", controllers.DownloadFiles).Methods("GET")

	log.Infoln("load api routers success")
	// // static files
	// router.Handle("/", http.FileServer(http.Dir("../frontend/")))
	// router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("../frontend/static/"))))
	return
}
