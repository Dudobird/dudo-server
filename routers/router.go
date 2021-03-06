package routers

import (
	"fmt"
	"net/http"

	"github.com/Dudobird/dudo-server/controllers"
	"github.com/Dudobird/dudo-server/utils"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func notFound(w http.ResponseWriter, r *http.Request) {
	utils.JSONMessageWithData(w, http.StatusNotFound, "url not exist", nil)
	return
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
	router.Use(jwtAuthenticationMiddleware)
	router.Use(appBindMiddleware)
	router.HandleFunc("/api/auth/signup", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/api/auth/signin", controllers.Login).Methods("POST")
	router.HandleFunc("/api/auth/logout", controllers.Logout).Methods("GET")
	router.HandleFunc("/api/auth/password", controllers.UpdatePassword).Methods("POST")

	router.HandleFunc("/api/folders", controllers.CreateFolder).Methods("POST")
	router.HandleFunc("/api/folders/{id}", controllers.ListFolderFiles).Methods("GET")

	router.HandleFunc("/api/files/{id}", controllers.GetFileInfo).Methods("GET")
	router.HandleFunc("/api/files/{id}", controllers.UpdateFileInfo).Methods("PUT")
	// delete files
	router.HandleFunc("/api/files/{id}", controllers.DeleteFiles).Methods("DELETE")
	// for top level becouse no folder just set it to `root`
	router.HandleFunc("/api/upload/files/{folderID}", controllers.UploadFiles).Methods("POST")
	router.HandleFunc("/api/download/files/{id}", controllers.DownloadFiles).Methods("GET")

	router.HandleFunc("/api/profile", controllers.GetProfile).Methods("GET")
	router.HandleFunc("/api/profile", controllers.UpdateProfile).Methods("PUT")

	router.HandleFunc("/api/shares", controllers.GetShareFiles).Methods("GET")
	router.HandleFunc("/api/shares", controllers.CreateShareFile).Methods("POST")
	router.HandleFunc("/api/share/{id}", controllers.DeleteShareFile).Methods("DELETE")
	router.HandleFunc("/shares", controllers.GetShareFileFromToken).Methods("GET")

	adminRouter := router.PathPrefix("/api/admin").Subrouter()
	adminRouter.Use(adminMiddleware)
	adminRouter.HandleFunc("/users", controllers.AdminGetUsers).Methods("GET")
	adminRouter.HandleFunc("/users/{id}", controllers.AdminDeleteUser).Methods("DELETE")
	adminRouter.HandleFunc("/users/{id}/limit", controllers.AdminChangeUserStorageLimit).Methods("PUT")
	adminRouter.HandleFunc("/users/{id}/password", controllers.AdminChangeUserPassword).Methods("PUT")
	// router.HandleFunc("/api/admin/shares", controllers.GetAdminShares).Methods("GET")
	// router.HandleFunc("/api/admin/files", controllers.GetAdminFiles).Methods("GET")

	router.HandleFunc("/api/search/files", controllers.HandleSearchFiles).Methods("POST")

	router.NotFoundHandler = http.HandlerFunc(notFound)
	log.Infoln("load api routers success")
	return
}
