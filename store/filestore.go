package store

import (
	"encoding/base64"
	"path/filepath"

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

// FileStore with a db handler
type FileStore struct {
	DB     *gorm.DB
	userID string
}

//NewFileStore return a FileStore instance
func NewFileStore(userID string) *FileStore {
	return &FileStore{
		DB:     models.GetDB(),
		userID: userID,
	}
}

// GetOrCreateFolder implement StorageHandler interface
// return the folder id of parent and path combination
// if path not exist create it and return id
func (store *FileStore) GetOrCreateFolder(parent, path string) (string, error) {
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

func (store *FileStore) createFoldersUnderParentID(parentID string, folders []string) (string, error) {
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

// StorageFileExistUnderFolderID  return true when file exist or false if not exist
func (store *FileStore) StorageFileExistUnderFolderID(folderID, fileName string) bool {
	existCheckStorage := &models.StorageFile{}
	notFoundChecker := store.DB.Where(
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

// StorageFileExistCheck  return true when file exist or false if not exist
func (store *FileStore) StorageFileExistCheck(fileID string) bool {
	existCheckStorage := &models.StorageFile{}
	notFoundChecker := store.DB.Model(&models.StorageFile{}).Where("id = ?", fileID).First(&existCheckStorage).RecordNotFound()
	if notFoundChecker == false {
		return true
	}
	return false
}

// SaveStorage save the data and update profile usage size
func (store *FileStore) SaveStorage(storage *models.StorageFile) error {
	err := store.DB.Save(storage).Error
	if err != nil {
		log.Errorf("save file error: %s", err)
		return err
	}
	profile := models.Profile{UserID: storage.UserID}
	err = store.DB.Model(&profile).UpdateColumn("usage_disk_size", gorm.Expr("usage_disk_size + ?", storage.FileSize)).Error
	if err != nil {
		log.Errorf("update user profile disk usage fail:%s", err)
		return err
	}
	return nil
}

// DeleteFolders delete folder and update the disk usage
func (store *FileStore) DeleteFolders(parentID string) ([]models.StorageFile, error) {
	pendingDeleteFiles := []models.StorageFile{}
	deleteFiles := []models.StorageFile{}
	err := store.DB.Where("id=?", parentID).Or("folder_id=?", parentID).Find(&pendingDeleteFiles).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return pendingDeleteFiles, nil
		}
		return pendingDeleteFiles, err
	}
	// delete in database
	// make sure all file belone to this user
	for _, f := range pendingDeleteFiles {
		if f.UserID != store.userID {
			continue
		}
		if f.IsDir == false {
			deleteFiles = append(deleteFiles, f)
		}
		if f.ID != parentID {
			// subfolder
			files, err := store.DeleteFolders(f.ID)
			if err != nil {
				log.Errorf("delete files fail: %s", err)
			} else {
				deleteFiles = append(deleteFiles, files...)
			}
		}
		// delete direct only in database
		store.DB.Unscoped().Delete(f)
	}
	return deleteFiles, nil
}

// GetAllFiles query all files and return a map info for store the path info and file into
func (store *FileStore) GetAllFiles(parentID string, parentName string) (map[string][]models.StorageFile, error) {
	queryFiles := []models.StorageFile{}
	files := make(map[string][]models.StorageFile)
	currentFolder := parentName
	err := store.DB.Model(&models.StorageFile{}).Where(
		"user_id = ? and folder_id = ?",
		store.userID,
		parentID,
	).Find(&queryFiles).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return files, nil
		}
	}
	// bug is here
	for _, f := range queryFiles {
		files[currentFolder] = append(files[currentFolder], f)
		if f.IsDir == false {
			continue
		}
		info, err := store.GetAllFiles(f.ID, filepath.Join(currentFolder, f.FileName))
		if err != nil {
			log.Errorf("download folder error:%s", err)
			continue
		}
		for k, v := range info {
			files[k] = v
		}
	}

	return files, nil
}
