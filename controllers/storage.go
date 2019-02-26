package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Dudobird/dudo-server/auth"
	"github.com/Dudobird/dudo-server/models"
	"github.com/Dudobird/dudo-server/utils"
)

// CreateFiles create a folder or file
func CreateFiles(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.TokenContextKey).(uint)
	_, errWithCode := models.GetUser(userID)
	if errWithCode != nil {
		utils.JSONRespnseWithErr(w, errWithCode)
		return
	}
	file := &models.StorageFile{UserID: userID}
	err := json.NewDecoder(r.Body).Decode(file)
	if err != nil {
		utils.JSONRespnseWithErr(w, &utils.ErrPostDataNotCorrect)
		return
	}
	errWithCode = file.CreateFile(userID)
	if errWithCode != nil {
		utils.JSONRespnseWithErr(w, errWithCode)
		return
	}
	utils.JSONRespnseWithTextMessage(w, http.StatusCreated, "")
	return
}

// GetTopLevelFiles list all top level files when user login success
func GetTopLevelFiles(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.TokenContextKey).(uint)
	user, errWithCode := models.GetUser(userID)
	if errWithCode != nil {
		utils.JSONRespnseWithErr(w, errWithCode)
		return
	}

	swu := models.StorageFilesWithUser{
		Owner: user,
	}

	data, errWithCode := swu.GetTopFiles()
	if errWithCode != nil {
		utils.JSONRespnseWithErr(w, errWithCode)
		return
	}
	utils.JSONMessageWithData(w, 200, "", data)
	return
}

// GetCurrentFile get information about one file or folder based the post id
func GetCurrentFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if utils.ValidateUUID(id) == false {
		utils.JSONRespnseWithErr(w, &utils.ErrPostDataNotCorrect)
		return
	}
	userID := r.Context().Value(auth.TokenContextKey).(uint)
	user, errWithCode := models.GetUser(userID)
	if errWithCode != nil {
		utils.JSONRespnseWithErr(w, errWithCode)
		return
	}
	swu := models.StorageFilesWithUser{
		Owner: user,
	}
	data, errWithCode := swu.ListCurrentFile(id)
	if errWithCode != nil {
		utils.JSONRespnseWithErr(w, errWithCode)
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
		utils.JSONRespnseWithErr(w, &utils.ErrPostDataNotCorrect)
		return
	}
	userID := r.Context().Value(auth.TokenContextKey).(uint)
	user, errWithCode := models.GetUser(userID)
	if errWithCode != nil {
		utils.JSONRespnseWithErr(w, errWithCode)
		return
	}
	swu := models.StorageFilesWithUser{
		Owner: user,
	}
	data, errWithCode := swu.ListChildren(id)
	if errWithCode != nil {
		utils.JSONRespnseWithErr(w, errWithCode)
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
		utils.JSONRespnseWithErr(w, &utils.ErrPostDataNotCorrect)
		return
	}

	userID := r.Context().Value(auth.TokenContextKey).(uint)
	user, errWithCode := models.GetUser(userID)
	if errWithCode != nil {
		utils.JSONRespnseWithErr(w, &utils.ErrPostDataNotCorrect)
		return
	}
	swu := models.StorageFilesWithUser{
		Owner: user,
	}
	errWithCode = swu.DeleteFileFromID(id)
	if errWithCode != nil {
		utils.JSONRespnseWithErr(w, &utils.ErrPostDataNotCorrect)
		return
	}
	utils.JSONMessageWithData(w, 200, "", nil)
	return
}
