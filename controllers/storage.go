package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Dudobird/dudo-server/auth"
	"github.com/Dudobird/dudo-server/core"
	"github.com/Dudobird/dudo-server/models"
	"github.com/Dudobird/dudo-server/utils"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// CreateFolder create a folder from post data
func CreateFolder(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.TokenContextKey).(string)
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
	errWithCode = file.CreateFolder(userID)
	if errWithCode != nil {
		utils.JSONRespnseWithErr(w, errWithCode)
		return
	}
	utils.JSONRespnseWithTextMessage(w, http.StatusCreated, "")
	return
}

// GetFileInfo get information about one file or folder based the post id
func GetFileInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	userID := r.Context().Value(auth.TokenContextKey).(string)
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

// ListFolderFiles list all top level files when user login success
func ListFolderFiles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	userID := r.Context().Value(auth.TokenContextKey).(string)
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
	app := core.GetApp()
	if app == nil {
		utils.JSONRespnseWithErr(w, &utils.ErrInternalServerError)
		return
	}
	userID := r.Context().Value(auth.TokenContextKey).(string)
	user, errWithCode := models.GetUser(userID)
	if errWithCode != nil {
		utils.JSONRespnseWithErr(w, &utils.ErrPostDataNotCorrect)
		return
	}
	swu := models.StorageFilesWithUser{
		Owner: user,
	}
	files, errWithCode := swu.DeleteFilesFromID(id)
	if errWithCode != nil {
		log.Error(errWithCode.Error())
		utils.JSONRespnseWithErr(w, &utils.ErrPostDataNotCorrect)
		return
	}
	messages := []string{}
	for _, file := range files {
		storeFileName := file.ID + "_" + file.FileName
		err := app.Storage.Delete(storeFileName, file.Bucket)
		if err != nil {
			log.Errorf("delete from storage error : %s", err)
			log.Errorf("delete detail info : %s %s", file.Bucket, storeFileName)
			messages = append(messages, fmt.Sprintf("%s:%s", file.FileName, err))
			continue
		}
		messages = append(messages, fmt.Sprintf("%s:success", file.FileName))
	}

	utils.JSONMessageWithData(w, 200, "", messages)
	return
}
