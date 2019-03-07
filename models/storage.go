package models

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/Dudobird/dudo-server/utils"
)

// StorageFile for store user files
type StorageFile struct {
	RawStorageFileInfo
	CreatedAt time.Time `gorm:"DEFAULT:current_timestamp"`
	UpdatedAt time.Time `gorm:"DEFAULT:current_timestamp"`
	DeletedAt *time.Time
	UserID    string `json:"user_id"`
}

// MarshalJSON custom json response
func (s *StorageFile) MarshalJSON() ([]byte, error) {
	type AliasStruct StorageFile
	return json.Marshal(&struct {
		FileSizeReadable string `json:"file_size_readable"`
		*AliasStruct
	}{
		FileSizeReadable: utils.GetReadableFileSize(float64(s.RawStorageFileInfo.FileSize)),
		AliasStruct:      (*AliasStruct)(s),
	})
}

//RawStorageFileInfo after upload success ,each file will generate one raw info
type RawStorageFileInfo struct {
	ID       string `json:"id" gorm:"primary_key"`
	FileName string `json:"file_name" gorm:"not null;index:idx_file_name"`
	// if you use local storage bucket will be folder name
	Bucket string `json:"bucket"`

	FileSize int64  `json:"file_size"`
	FolderID string `json:"folder_id" gorm:"not null;default:''"`
	IsDir    bool   `json:"is_dir" gorm:"not null"`
	Path     string `json:"path" gorm:"not null"`
}

func (s *StorageFile) validation() *utils.CustomError {
	if s.UserID == "" {
		return &utils.ErrPostDataNotCorrect
	}
	folder := &StorageFile{}
	// file name validation
	if s.FileName == "" || len(s.FileName) > 50 {
		return &utils.ErrPostDataNotCorrect
	}
	// if folderID not exist, it will create in root position
	// or validate if folder exist or not
	if s.FolderID != "root" && s.FolderID != "" {
		// create top level
		// check parent id exist or not
		err := GetDB().Model(&StorageFile{}).Where("id = ?", s.FolderID).First(folder).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return &utils.ErrResourceNotFound
			}
			return &utils.ErrInternalServerError
		}
	}
	// check if resource already exist in same folder with same filename
	err := GetDB().Where("file_name = ? AND folder_id = ?", s.FileName, s.FolderID).First(folder).Error
	if err == nil {
		return &utils.ErrResourceAlreadyExist
	}
	if err != gorm.ErrRecordNotFound {
		return &utils.ErrInternalServerError
	}
	return nil
}

// CreateFolder will create a new folder from post data
// if file type is folder just create in database
func (s *StorageFile) CreateFolder(uid string) *utils.CustomError {
	// valid user post data
	if err := s.validation(); err != nil {
		return err
	}
	if s.IsDir == false {
		return &utils.ErrPostDataNotCorrect
	}
	s.ID = utils.GenRandomID("folder", 15)
	err := GetDB().Model(&StorageFile{}).Create(s).Error
	if err != nil {
		return &utils.ErrInternalServerError
	}
	return nil
}

// StorageFilesWithUser for controller
type StorageFilesWithUser struct {
	Owner *User
	Files []*StorageFile
}

// ListCurrentFile list the file with id
func (swu *StorageFilesWithUser) ListCurrentFile(id string) (*StorageFile, *utils.CustomError) {
	file := StorageFile{}
	err := GetDB().Model(&StorageFile{}).Where("id=? and user_id=?", id, swu.Owner.ID).Find(&file).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &utils.ErrResourceNotFound
		}
		return nil, &utils.ErrInternalServerError
	}
	return &file, nil
}

// ListChildren list the file and all subfiles
func (swu *StorageFilesWithUser) ListChildren(folderID string) ([]StorageFile, *utils.CustomError) {
	files := []StorageFile{}
	err := GetDB().Model(&StorageFile{}).Where("folder_id=? and user_id=?", folderID, swu.Owner.ID).Find(&files).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &utils.ErrResourceNotFound
		}
		return nil, &utils.ErrInternalServerError
	}
	return files, nil
}

// DeleteFilesFromID delete files based on its id and delete all subfiles
// return the files infomation to delete in storage
func (swu *StorageFilesWithUser) DeleteFilesFromID(id string) ([]StorageFile, *utils.CustomError) {
	pendingDeleteFiles := []StorageFile{}
	err := GetDB().Where("id=?", id).Or("folder_id=?", id).Find(&pendingDeleteFiles).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return pendingDeleteFiles, &utils.ErrResourceNotFound
		}
		return pendingDeleteFiles, &utils.ErrInternalServerError
	}
	// delete in database
	// make sure all file belone to this user
	for i, f := range pendingDeleteFiles {
		if f.UserID != swu.Owner.ID {
			continue
		}
		// if it's file then send to storage manager to delete
		// if is directory then skip only delete in database
		if f.RawStorageFileInfo.IsDir == true {
			pendingDeleteFiles = append(pendingDeleteFiles[:i], pendingDeleteFiles[i+1:]...)
		}
		// delete direct only in database
		GetDB().Unscoped().Delete(f)
	}
	return pendingDeleteFiles, nil
}
