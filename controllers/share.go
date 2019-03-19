package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Dudobird/dudo-server/store"
	"github.com/Dudobird/dudo-server/utils"
	"github.com/gorilla/mux"
)

type shareFileInfo struct {
	FileID      string `json:"file_id"`
	ExpireDays  int    `json:"expire_days"`
	Description string `json:"description"`
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
	token, err := fileStore.CreateShareToken(sfi.FileID, sfi.ExpireDays, sfi.Description)
	if err != nil {
		utils.JSONRespnseWithErr(w, (err).(*utils.CustomError))
		return
	}
	utils.JSONMessageWithData(w, 201, "", struct {
		Token string `json:"token"`
	}{Token: token})
	return
}

// GetShareFileFromToken download share file for others
func GetShareFileFromToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		utils.JSONMessageWithData(w, 400, "token is empty", nil)
		return
	}
	fileStore := store.NewFileStore("")
	fileID, userID, err := fileStore.VerifyShareToken(token)
	if err != nil {
		utils.JSONRespnseWithErr(w, err)
		return
	}
	ctx := context.WithValue(r.Context(), utils.TokenContextKey, userID)
	r = r.WithContext(ctx)
	data := make(map[string]string)
	data["id"] = fileID
	r = mux.SetURLVars(r, data)
	DownloadFiles(w, r)
	return
}

// GetShareFiles get all shared folders with user id
func GetShareFiles(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.TokenContextKey).(string)
	fileStore := store.NewFileStore(userID)
	files, err := fileStore.GetAllSharedFiles()
	if err != nil {
		utils.JSONRespnseWithErr(w, &utils.ErrInternalServerError)
		return
	}
	utils.JSONMessageWithData(w, 200, "", files)
	return
}

// DeleteShareFile get the id from url and delete the share file reference
func DeleteShareFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	userID := r.Context().Value(utils.TokenContextKey).(string)

	fileStore := store.NewFileStore(userID)
	err := fileStore.DeleteShareFilesRef(id)
	if err != nil {
		utils.JSONRespnseWithErr(w, err)
		return
	}
	utils.JSONMessageWithData(w, 200, "delete success", nil)
	return

}
