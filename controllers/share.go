package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Dudobird/dudo-server/store"
	"github.com/Dudobird/dudo-server/utils"
)

type shareFileInfo struct {
	FileID     string `json:"file_id"`
	ExpireDays int    `json:"expire_days"`
}

// CreateShareFile create a new shared files
func CreateShareFile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.TokenContextKey).(string)

	sfi := &shareFileInfo{}
	err := json.NewDecoder(r.Body).Decode(sfi)
	if err != nil {
		utils.JSONRespnseWithErr(w, &utils.ErrPostDataNotCorrect)
		return
	}
	fileStore := store.NewFileStore(userID)
	token, err := fileStore.CreateShareToken(sfi.FileID, sfi.ExpireDays)
	if err != nil {
		utils.JSONRespnseWithErr(w, (err).(*utils.CustomError))
		return
	}
	utils.JSONMessageWithData(w, 201, "", struct {
		Token string `json:"token"`
	}{Token: token})
	return
}

// GetShareFiles get all shared folders with user id
func GetShareFiles(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.TokenContextKey).(string)
	fileStore := store.NewFileStore(userID)
	files, err := fileStore.GetAllSharedFiles()
	if err == nil {
		utils.JSONRespnseWithErr(w, &utils.ErrInternalServerError)
		return
	}
	utils.JSONMessageWithData(w, 201, "", files)
	return
}
