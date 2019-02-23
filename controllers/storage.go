package controllers

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Dudobird/dudo-server/auth"
	"github.com/Dudobird/dudo-server/models"
	"github.com/Dudobird/dudo-server/utils"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// GetTopLevelFiles list all top level files when user login success
func GetTopLevelFiles(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.TokenContextKey).(uint)
	user := models.GetUser(userID)

	if user == nil {
		utils.JSONRespnseWithTextMessage(w, http.StatusNotFound, "user not found")
		return
	}

	swu := models.StorageFilesWithUser{
		Owner: user,
	}
	err := swu.GetTopFiles()
	if err != nil {
		log.Errorln(err)
		utils.JSONRespnseWithTextMessage(w, http.StatusServiceUnavailable, "request service not avaliable now")
		return
	}
	utils.JSONMessageWithData(w, 200, "", swu.Files)
	return
}

// ListSubFiles list all top level files when user login success
func ListSubFiles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if utils.ValidateUUID(id) == false {
		utils.JSONRespnseWithTextMessage(w, http.StatusBadRequest, "id not given")
		return
	}
	userID := r.Context().Value(auth.TokenContextKey).(uint)
	user := models.GetUser(userID)
	if user == nil {
		utils.JSONRespnseWithTextMessage(w, http.StatusNotFound, "user not found")
		return
	}
	swu := models.StorageFilesWithUser{
		Owner: user,
	}
	data, err := swu.ListChildren(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.JSONRespnseWithTextMessage(w, http.StatusNotFound, "request object not found")
			return
		}
		log.Errorln(err)
		utils.JSONRespnseWithTextMessage(w, http.StatusServiceUnavailable, "request service not avaliable now")
		return
	}
	utils.JSONMessageWithData(w, 200, "", data)
	return
}

// DeleteFiles delete current file and all reference files
func DeleteFiles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if utils.ValidateUUID(id) == false {
		utils.JSONRespnseWithTextMessage(w, http.StatusBadRequest, "id not given")
		return
	}

	userID := r.Context().Value(auth.TokenContextKey).(uint)
	user := models.GetUser(userID)
	if user == nil {
		utils.JSONRespnseWithTextMessage(w, http.StatusNotFound, "user not found")
		return
	}
	swu := models.StorageFilesWithUser{
		Owner: user,
	}
	err := swu.DeleteFileFrom(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.JSONRespnseWithTextMessage(w, http.StatusNotFound, "request object not found")
			return
		}
		log.Errorln(err)
		utils.JSONRespnseWithTextMessage(w, http.StatusServiceUnavailable, "request service not avaliable now")
		return
	}
	utils.JSONMessageWithData(w, 200, "", nil)
	return
}
