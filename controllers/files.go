package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/jinzhu/gorm"

	"github.com/Dudobird/dudo-server/auth"
	"github.com/Dudobird/dudo-server/models"
	"github.com/Dudobird/dudo-server/store"
	"github.com/gorilla/mux"

	"github.com/Dudobird/dudo-server/core"
	"github.com/Dudobird/dudo-server/utils"

	log "github.com/sirupsen/logrus"
)

// UploadFiles receive user upload file
// save it to temp folder and wait for upload to storage
func UploadFiles(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.TokenContextKey).(string)
	vars := mux.Vars(r)
	folderID := vars["folderID"]
	filePath := r.Header.Get("X-FilePath")
	store := store.NewStore(userID)
	folderID, err := store.GetOrCreateFolder(folderID, filePath)
	if err != nil {
		utils.JSONRespnseWithErr(w, err.(*utils.CustomError))
		return
	}
	app := core.GetApp()

	r.ParseMultipartForm(64 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		log.Errorf("upload file fail : %s ", err)
		utils.JSONRespnseWithErr(w, &utils.ErrPostDataNotCorrect)
		return
	}
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		log.Errorf("upload file fail : %s ", err)
		utils.JSONRespnseWithErr(w, &utils.ErrPostDataNotCorrect)
		return
	}
	// https://golang.org/pkg/net/http/#DetectContentType
	mimeType := http.DetectContentType(buff)
	defer file.Close()

	exist := store.StorageFileExistCheck(folderID, handler.Filename)
	if exist == true {
		utils.JSONRespnseWithErr(w, &utils.ErrResourceAlreadyExist)
		return
	}

	// the tempfile will be userid_timestamp_realfilename
	id := utils.GenRandomID("file", 15)
	// bucket name has some restrict
	// https://docs.aws.amazon.com/AmazonS3/latest/dev/BucketRestrictions.html
	bucketName := fmt.Sprintf(
		"%s-%s",
		app.Config.Application.BucketPrefix,
		strings.ToLower(strings.TrimLeft(userID, "user_")),
	)
	fileName := fmt.Sprintf("%s_%s", id, handler.Filename)
	tempFileName := app.FullTempFolder + string(filepath.Separator) + fileName
	f, err := os.OpenFile(tempFileName, os.O_WRONLY|os.O_CREATE, 0666)
	defer func() {
		f.Close()
		os.Remove(tempFileName)
	}()
	if err != nil {
		log.Errorf("save temp file fail : %s", err)
		utils.JSONRespnseWithErr(w, &utils.ErrInternalServerError)
		return
	}
	file.Seek(0, 0)
	size, _ := io.Copy(f, file)
	defer f.Close()
	path, err := app.Storage.Upload(tempFileName, fileName, bucketName)
	if err != nil {
		log.Errorf("upload to storage fail : %s", err)
		utils.JSONRespnseWithErr(w, &utils.ErrInternalServerError)
		return
	}
	// save to database
	s := models.StorageFile{
		UserID: userID,
		RawStorageFileInfo: models.RawStorageFileInfo{
			ID:       id,
			FileName: handler.Filename,
			Bucket:   bucketName,
			IsDir:    false,
			MIMEType: mimeType,
			FileType: utils.GetFileExtention(handler.Filename),
			FileSize: size,
			FolderID: folderID,
			Path:     path,
		},
	}
	app.DB.Save(&s)
	utils.JSONMessageWithData(w, 201, "", id)
	return
}

// DownloadFiles will down load files from storages
func DownloadFiles(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.TokenContextKey).(string)
	vars := mux.Vars(r)
	id := vars["id"]
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

	tempDownloadFilePath := app.FullTempFolder + string(filepath.Separator) + fileMeta.FileName
	storeFileName := fileMeta.ID + "_" + fileMeta.FileName
	// Right now only allow for download file
	if fileMeta.IsDir == true {
		// zip folder
		// app.Storage.DownloadFolders(tempDownloadFilePath, storeFileName, fileMeta.Bucket)
		store := store.NewStore(userID)
		files, err := store.GetAllFiles(fileMeta.ID, "")
		if err != nil {
			log.Errorf("download folder error: %s", err)
			utils.JSONRespnseWithErr(w, &utils.ErrInternalServerError)
			return
		}
		utils.JSONMessageWithData(w, 200, "", files)
		return
	}

	err = app.Storage.Download(tempDownloadFilePath, storeFileName, fileMeta.Bucket)
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
