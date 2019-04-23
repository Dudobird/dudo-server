package store

import (
	"database/sql"
	"encoding/base64"
	"fmt"
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
		store.DB.Unscoped().Delete(f)
	}
	return deleteFiles, nil
}

//SearchResults search result return
type SearchResults struct {
	ParentID       string             `json:"parent_id"`
	ParentFileName string             `json:"parent_filename"`
	File           models.StorageFile `json:"file"`
}

// SearchFiles search files from metadata
func (store *FileStore) SearchFiles(search string) ([]SearchResults, error) {
	searchFiles := []SearchResults{}

	// var parentFiles []*models.StorageFile
	// var storageFiles []*models.StorageFile
	queryString := `
	SELECT P.id AS parent_id, P.file_name AS parent_filename,
	C.id, C.file_name,C.mime_type,C.file_type,C.file_size,C.is_dir,C.created_at,C.updated_at,C.deleted_at FROM storage_files AS P RIGHT OUTER JOIN storage_files 
	AS C ON C.folder_id=P.id where C.user_id = ? and C.file_name LIKE ?`
	rows, err := store.DB.DB().Query(queryString, store.userID, fmt.Sprintf("%%%s%%", search))
	if err != nil {
		log.Error(err)
		return searchFiles, &utils.ErrInternalServerError
	}
	defer rows.Close()
	for rows.Next() {
		var pid, pfilename sql.NullString
		child := models.StorageFile{}
		if err := rows.Scan(&pid, &pfilename, &child.ID, &child.FileName,
			&child.MIMEType,
			&child.FileType,
			&child.FileSize,
			&child.IsDir,
			&child.CreatedAt, &child.UpdatedAt, &child.DeletedAt); err != nil {
			log.Error(err)
		}
		searchFiles = append(searchFiles, SearchResults{
			ParentID:       pid.String,
			ParentFileName: pfilename.String,
			File:           child,
		})
	}
	return searchFiles, nil
}

func (store *FileStore) validationFileName(fileName string) error {
	if fileName == "" || len(fileName) > 100 {
		return &utils.ErrPostDataNotCorrect
	}
	return nil
}

//RenameFileName rename filename with id
func (store *FileStore) RenameFileName(id string, name string) (*models.StorageFile, error) {
	// valid user post data
	if err := store.validationFileName(name); err != nil {
		return nil, err
	}
	file := &models.StorageFile{}
	err := store.DB.Model(&models.StorageFile{}).Where(
		"id = ? and user_id = ?",
		id,
		store.userID,
	).First(file).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &utils.ErrResourceNotFound
		}
		return nil, &utils.ErrInternalServerError
	}
	// query if already have some file or folder with same name
	temp := &models.StorageFile{}
	err = store.DB.Where("file_name = ? AND folder_id = ?", name, file.FolderID).First(temp).Error
	if err != gorm.ErrRecordNotFound {
		if err != nil {
			return nil, &utils.ErrInternalServerError
		}
		return nil, &utils.ErrResourceAlreadyExist
	}
	file.FileName = name
	err = store.DB.Save(file).Error
	if err != nil {
		return nil, &utils.ErrInternalServerError
	}
	return file, nil
}

// GetAllFiles query all files and return a map info for store the path info and file into
func (store *FileStore) GetAllFiles(parentID string, parentName string) (map[string][]models.StorageFile, bool, error) {
	queryFiles := []models.StorageFile{}
	hasFiles := false
	files := make(map[string][]models.StorageFile)
	currentFolder := parentName
	err := store.DB.Model(&models.StorageFile{}).Where(
		"user_id = ? and folder_id = ?",
		store.userID,
		parentID,
	).Find(&queryFiles).Error
	if err != nil {
		log.Errorf("download folder error:%s", err)
		return files, hasFiles, err
	}

	if len(queryFiles) == 0 {
		return files, hasFiles, nil
	}
	for _, f := range queryFiles {
		files[currentFolder] = append(files[currentFolder], f)
		if f.IsDir == false {
			hasFiles = true
			continue
		}
		info, subFolderHasFiles, err := store.GetAllFiles(f.ID, filepath.Join(currentFolder, f.FileName))
		hasFiles = hasFiles || subFolderHasFiles
		if err != nil {
			log.Errorf("download folder error:%s", err)
			continue
		}
		for k, v := range info {
			files[k] = v
		}
	}

	return files, hasFiles, nil
}
