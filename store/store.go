package store

import (
	"encoding/base64"

	"github.com/Dudobird/dudo-server/models"
	"github.com/Dudobird/dudo-server/utils"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// StorageHandler is a interface for user
// do some sql for Storage
type StorageHandler interface {
	GetOrCreateFolder(parent, path string, override bool) (string, error)
}

// Store with a db handler
type Store struct {
	DB     *gorm.DB
	userID string
}

//NewStore return a Store instance
func NewStore(userID string) *Store {
	return &Store{
		DB:     models.GetDB(),
		userID: userID,
	}
}

// GetOrCreateFolder implement StorageHandler interface
// return the folder id of parent and path combination
// if path not exist create it and return id
func (store *Store) GetOrCreateFolder(parent, path string) (string, error) {
	// split path to slice
	pathDecodeBase64, err := base64.StdEncoding.DecodeString(path)
	if err != nil {
		return "", &utils.ErrPostDataNotCorrect
	}
	pathDecodeBase64Str := string(pathDecodeBase64)
	if parent != "root" {
		var exist int
		err := store.DB.Model(&models.StorageFile{}).Where("id = ?", parent).Count(&exist).Error
		if err != nil {
			return "", &utils.ErrInternalServerError
		}
		if exist == 0 {
			return "", &utils.ErrResourceNotFound
		}
	}
	folders := utils.GetFilePathFolderList(pathDecodeBase64Str)
	if pathDecodeBase64Str == "" || len(folders) == 0 {
		return parent, nil
	}
	// parent is root or exist ID
	folderID, err := store.createFoldersUnderParentID(parent, folders)
	return folderID, err
}

func (store *Store) createFoldersUnderParentID(parentID string, folders []string) (string, error) {
	var parent = parentID
	var currentFolderID string
	var err error
	checkIfExist := true

	for index, folder := range folders {
		if checkIfExist {
			s := &models.StorageFile{}
			err = store.DB.Model(&models.StorageFile{}).Where(
				"user_id = ? and folder_id = ? and file_name = ?",
				store.userID,
				parent,
				folders[index],
			).First(s).Error

			if err != nil && err != gorm.ErrRecordNotFound {
				return "", &utils.ErrInternalServerError
			}
			if err == nil && s.IsDir == false {
				// some file has this folder name
				return "", &utils.ErrResourceAlreadyExist
			}
			if err == nil && s.IsDir == true {
				log.Infof("get folder :%v", s.ID)
				parent = s.ID
				continue
			}
		}
		checkIfExist = false
		currentFolderID = utils.GenRandomID("folder", 15)
		err := store.DB.Model(&models.StorageFile{}).Save(&models.StorageFile{
			UserID: store.userID,
			RawStorageFileInfo: models.RawStorageFileInfo{
				ID:       currentFolderID,
				FileName: folder,
				IsDir:    true,
				FolderID: parent,
				Path:     "",
			},
		}).Error
		if err != nil {
			log.Errorf("create folder fail:%s", err)
			return "", &utils.ErrInternalServerError
		}
		parent = currentFolderID
		continue
	}
	return parent, nil
}

// StorageFileExistCheck  return true when file exist or false if not exist
func (store *Store) StorageFileExistCheck(folderID, fileName string) bool {
	existCheckStorage := &models.StorageFile{}
	notFoundChecker := models.GetDB().Where(
		&models.StorageFile{
			RawStorageFileInfo: models.RawStorageFileInfo{
				FolderID: folderID,
				FileName: fileName,
			},
			UserID: store.userID,
		},
	).First(&existCheckStorage).RecordNotFound()
	if notFoundChecker == false {
		return true
	}
	return false
}
