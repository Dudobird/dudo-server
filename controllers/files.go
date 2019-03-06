package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jinzhu/gorm"

	"github.com/Dudobird/dudo-server/auth"
	"github.com/Dudobird/dudo-server/models"
	"github.com/gorilla/mux"

	uuid "github.com/satori/go.uuid"

	"github.com/Dudobird/dudo-server/core"
	"github.com/Dudobird/dudo-server/utils"

	log "github.com/sirupsen/logrus"
)

// UploadFiles receive user upload file
// save it to temp folder and wait for upload to storage
// user post data to /api/upload/root or /api/upload/**-**-**-**(parentID)
func UploadFiles(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.TokenContextKey).(uint)
	vars := mux.Vars(r)

	parentID := vars["parentID"]
	if parentID == "root" {
		parentID = ""
	}
	if parentID != "" && utils.ValidateUUID(parentID) == false {
		utils.JSONRespnseWithErr(w, &utils.ErrPostDataNotCorrect)
		return
	}
	app := core.GetApp()
	if app == nil {
		utils.JSONRespnseWithErr(w, &utils.ErrInternalServerError)
		return
	}

	// 64 Mb in minio it will upload in one file and not split to many trunk
	r.ParseMultipartForm(64 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		log.Errorf("Upload file fail : %s ", err)
		utils.JSONRespnseWithErr(w, &utils.ErrPostDataNotCorrect)
		return
	}
	defer file.Close()
	// the tempfile will be userid_timestamp_realfilename
	id := uuid.NewV4()
	bucketName := fmt.Sprintf(
		"%s-%d",
		app.Config.Application.BucketPrefix,
		userID,
	)
	fileName := fmt.Sprintf("%s_%s", id, handler.Filename)
	tempFileName := app.FullTempFolder + string(filepath.Separator) + fileName
	f, err := os.OpenFile(tempFileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Errorf("save temp file fail : %s", err)
		utils.JSONRespnseWithErr(w, &utils.ErrInternalServerError)
		return
	}
	size, _ := io.Copy(f, file)
	f.Close()
	path, err := app.StorageHandler.Upload(tempFileName, fileName, bucketName)
	if err != nil {
		log.Errorf("upload to storage fail : %s", err)
		utils.JSONRespnseWithErr(w, &utils.ErrInternalServerError)
		return
	}
	// save to database
	s := models.StorageFile{
		UserID: userID,
		RawStorageFileInfo: models.RawStorageFileInfo{
			ID:       id.String(),
			FileName: handler.Filename,
			Bucket:   bucketName,
			IsDir:    false,
			FileSize: size,
			ParentID: parentID,
			Path:     path,
		},
	}
	app.DB.Save(&s)
	os.Remove(tempFileName)
	utils.JSONMessageWithData(w, 201, "", id.String())
	return
}

// DownloadFiles will down load files from storages
func DownloadFiles(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.TokenContextKey).(uint)
	vars := mux.Vars(r)
	id := vars["id"]
	if utils.ValidateUUID(id) == false {
		utils.JSONRespnseWithErr(w, &utils.ErrPostDataNotCorrect)
		return
	}
	app := core.GetApp()
	if app == nil {
		utils.JSONRespnseWithErr(w, &utils.ErrInternalServerError)
		return
	}
	fileMeta := &models.StorageFile{}
	err := app.DB.Model(&models.StorageFile{}).Where("id = ? and user_id = ?", id, userID).First(fileMeta).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.JSONRespnseWithErr(w, &utils.ErrResourceNotFound)
			return
		}
		log.Errorf("query storage info err: %s", err)
		utils.JSONRespnseWithErr(w, &utils.ErrInternalServerError)
		return
	}
	// Right now only allow for download file
	if fileMeta.IsDir == true {
		utils.JSONRespnseWithErr(w, &utils.ErrPostDataNotCorrect)
		return
	}
	tempDownloadFilePath := app.FullTempFolder + string(filepath.Separator) + fileMeta.FileName
	storeFileName := fileMeta.ID + "_" + fileMeta.FileName
	err = app.StorageHandler.Download(tempDownloadFilePath, storeFileName, fileMeta.Bucket)
	if err != nil {
		log.Errorf("down load file from storage err: %s", err)
		log.Errorf("filename = %s, bucket = %s", storeFileName, fileMeta.Bucket)
		utils.JSONRespnseWithErr(w, &utils.ErrInternalServerError)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileMeta.FileName+"\"")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	f, err := os.Open(tempDownloadFilePath)
	if err != nil {
		log.Errorf("open temp file err: %s", err)
		utils.JSONRespnseWithErr(w, &utils.ErrInternalServerError)
		return
	}
	defer func() {
		f.Close()
		os.Remove(tempDownloadFilePath)
	}()
	io.Copy(w, f)

	return
}
