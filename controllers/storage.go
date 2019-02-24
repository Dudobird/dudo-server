package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Dudobird/dudo-server/auth"
	"github.com/Dudobird/dudo-server/models"
	"github.com/Dudobird/dudo-server/utils"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// CreateFiles create a folder or file
func CreateFiles(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.TokenContextKey).(uint)
	user := models.GetUser(userID)
	if user == nil {
		utils.JSONRespnseWithTextMessage(w, http.StatusNotFound, "user not found")
		return
	}
	file := &models.StorageFile{}
	err := json.NewDecoder(r.Body).Decode(file)
	if err != nil {
		utils.JSONRespnseWithTextMessage(w, http.StatusBadRequest, "request data invalid")
		return
	}
	errWithCode := file.CreateFile(userID)
	if errWithCode != nil {
		utils.JSONRespnseWithTextMessage(w, errWithCode.Code(), errWithCode.Error())
		return
	}
	utils.JSONRespnseWithTextMessage(w, http.StatusCreated, "")
	return
}

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

// GetCurrentFile get information about one file or folder based the post id
func GetCurrentFile(w http.ResponseWriter, r *http.Request) {
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
	data, err := swu.ListCurrentFile(id)
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
	err := swu.DeleteFileFromID(id)
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
