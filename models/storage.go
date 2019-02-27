package models

import (
	"fmt"
	"time"

	"github.com/Dudobird/dudo-server/config"

	"github.com/jinzhu/gorm"

	"github.com/Dudobird/dudo-server/utils"

	uuid "github.com/satori/go.uuid"
)

// StorageFile for store user files
type StorageFile struct {
	CreatedAt time.Time `gorm:"DEFAULT:current_timestamp"`
	UpdatedAt time.Time `gorm:"DEFAULT:current_timestamp"`
	DeletedAt *time.Time
	RawStorageFileInfo
	UserID uint `json:"user_id"`
}

//RawStorageFileInfo after upload success ,each file will generate one raw info
type RawStorageFileInfo struct {
	ID       string `json:"id" gorm:"primary_key"`
	FileName string `json:"file_name" gorm:"not null;index:idx_file_name"`
	// if you use local storage bucket will be folder name
	Bucket string `json:"bucket"`
	// ParentID logic parent id
	ParentID string `json:"parent_id" gorm:"not null;default:''"`
	IsDir    bool   `json:"is_dir" gorm:"not null"`
	// remote minio storage path
	Path string `json:"path" gorm:"not null"`
}

func (sf *StorageFile) validation() *utils.CustomError {
	if sf.UserID == 0 {
		return &utils.ErrPostDataNotCorrect
	}
	parent := &StorageFile{}
	// file name validation
	if sf.FileName == "" || len(sf.FileName) > 50 {
		return &utils.ErrPostDataNotCorrect
	}
	// if parentid not exist, it will create in root position
	// or validate if parentid exist or not
	if sf.ParentID != "" {
		// create top level
		// check parent id exist or not
		err := GetDB().Model(&StorageFile{}).Where("id = ?", sf.ParentID).First(parent).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return &utils.ErrResourceNotFound
			}
			return &utils.ErrInternalServerError
		}
	}
	// check if resource already exist in same folder with same filename
	err := GetDB().Where("file_name = ? AND parent_id = ?", sf.FileName, sf.ParentID).First(parent).Error
	if err == nil {
		return &utils.ErrResourceAlreadyExist
	}
	if err != gorm.ErrRecordNotFound {
		return &utils.ErrInternalServerError
	}
	return nil
}

// CreateFile will save a new file from post data
// if file type is folder just create in database
func (sf *StorageFile) CreateFile(uid uint) *utils.CustomError {
	// valid user post data
	if err := sf.validation(); err != nil {
		return err
	}
	if sf.IsDir == false {
		// err := sf.createNewFile()
		// if err != nil {
		// 	return err
		// }
	}
	id := uuid.NewV4()
	sf.ID = id.String()
	err := GetDB().Model(&StorageFile{}).Create(sf).Error
	if err != nil {
		return &utils.ErrInternalServerError
	}
	return nil
}

func (sf *StorageFile) createNewFile() error {
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
func (swu *StorageFilesWithUser) ListChildren(parentID string) ([]StorageFile, *utils.CustomError) {
	files := []StorageFile{}
	err := GetDB().Model(&StorageFile{}).Where("parent_id=? and user_id=?", parentID, swu.Owner.ID).Find(&files).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &utils.ErrResourceNotFound
		}
		return nil, &utils.ErrInternalServerError
	}
	return files, nil
}

// SaveFromRawFiles files from raw info
func (swu *StorageFilesWithUser) SaveFromRawFiles(files []RawStorageFileInfo) error {
	config := config.GetConfig()
	bucketName := fmt.Sprintf("%s-%d", config.Application.BucketPrefix, swu.Owner.ID)
	for _, file := range files {
		file.Bucket = bucketName
		swu.Files = append(
			swu.Files,
			&StorageFile{
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
				RawStorageFileInfo: file,
			},
		)
	}
	err := GetDB().Model(swu.Owner).Association("Files").Append(swu.Files).Error
	return err
}

// GetTopFiles get user first level of files
func (swu *StorageFilesWithUser) GetTopFiles() ([]StorageFile, *utils.CustomError) {
	files := []StorageFile{}
	err := GetDB().Model(&StorageFile{}).Where("parent_id = ? and user_id = ?", "", swu.Owner.ID).Find(&files).Error
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
	err := GetDB().Where("id=?", id).Or("parent_id=?", id).Find(&pendingDeleteFiles).Error
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
