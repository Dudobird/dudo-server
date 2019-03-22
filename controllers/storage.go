package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Dudobird/dudo-server/store"

	"github.com/Dudobird/dudo-server/core"
	"github.com/Dudobird/dudo-server/models"
	"github.com/Dudobird/dudo-server/utils"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// CreateFolder create a folder from post data
func CreateFolder(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.TokenContextKey).(string)
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

//UpdateFileInfo will rename file name and change update
func UpdateFileInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	userID := r.Context().Value(utils.TokenContextKey).(string)

	fileInfo := &struct {
		Name string `json:"file_name"`
	}{}
	err := json.NewDecoder(r.Body).Decode(fileInfo)
	if err != nil {
		log.Error(err)
		utils.JSONRespnseWithErr(w, &utils.ErrPostDataNotCorrect)
		return
	}
	fileStore := store.NewFileStore(userID)
	data, errWithCode := fileStore.RenameFileName(id, fileInfo.Name)
	if errWithCode != nil {
		utils.JSONRespnseWithErr(w, errWithCode)
		return
	}
	utils.JSONMessageWithData(w, 200, "", data)
	return
}

// GetFileInfo get information about one file or folder based the post id
func GetFileInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	userID := r.Context().Value(utils.TokenContextKey).(string)
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
	userID := r.Context().Value(utils.TokenContextKey).(string)
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

// DeleteFiles delete current file or folder( all reference files)
func DeleteFiles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	app := core.GetApp()
	if app == nil {
		utils.JSONRespnseWithErr(w, &utils.ErrInternalServerError)
		return
	}
	userID := r.Context().Value(utils.TokenContextKey).(string)

	fileStore := store.NewFileStore(userID)
	// get files for storage delete
	files, err := fileStore.DeleteFolders(id)
	if err != nil {
		log.Errorf("delete folders fail : %s", err)
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

type searchInfo struct {
	Search string `json:"search"`
}

// HandleSearchFiles receive post data for files search  and response results
func HandleSearchFiles(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.TokenContextKey).(string)
	searchBody := &searchInfo{}
	err := json.NewDecoder(r.Body).Decode(searchBody)
	if err != nil || searchBody.Search == "" {
		utils.JSONRespnseWithErr(w, &utils.ErrPostDataNotCorrect)
		return
	}
	fileStore := store.NewFileStore(userID)
	files, err := fileStore.SearchFiles(searchBody.Search)
	if err != nil {
		log.Errorf("search file err:%s", err)
		utils.JSONRespnseWithErr(w, &utils.ErrInternalServerError)
		return
	}
	utils.JSONMessageWithData(w, 200, "", files)
	return
}
