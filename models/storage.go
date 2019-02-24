package models

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/Dudobird/dudo-server/utils"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
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
	ID         string `json:"id" gorm:"primary_key"`
	FileName   string `json:"file_name" gorm:"not null;index:idx_file_name"`
	Bucket     string `json:"bucket"`
	ParentID   string `json:"parent_id" gorm:"not null"`
	IsDir      bool   `json:"is_dir" gorm:"not null"`
	IsTopLevel bool   `json:"is_top_level" gorm:"not null"`
	// remote minio storage path
	Path string `json:"path" gorm:"not null"`
}

func (sf *StorageFile) validation() *utils.CustomError {
	if sf.IsTopLevel && sf.ParentID != "" {
		return &utils.ErrPostDataNotCorrect
	}
	// check parent id exist or not
	parent := &StorageFile{}
	err := GetDB().Model(&StorageFile{}).Where("parent_id = ?", sf.ParentID).First(parent).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &utils.ErrResourceNotFound
		}
		return &utils.ErrInternalServerError
	}

	// check if resource already exist in same folder
	err = GetDB().Where("file_name = ? AND parent_id = ?", sf.FileName, sf.ParentID).First(parent).Error
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
	id, err := uuid.NewV4()
	if err != nil {
		log.Errorf("create uuid with error: %s", err)
		return &utils.ErrInternalServerError
	}
	sf.ID = id.String()
	err = GetDB().Model(&StorageFile{}).Create(sf).Error
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
	for _, file := range files {
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
	err := GetDB().Model(&StorageFile{}).Where("is_top_level = ? and user_id = ?", true, swu.Owner.ID).Find(&files).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &utils.ErrResourceNotFound
		}
		return nil, &utils.ErrInternalServerError
	}
	return files, nil
}

// DeleteFileFromID delete file based on its id and delete all subfiles
func (swu *StorageFilesWithUser) DeleteFileFromID(id string) *utils.CustomError {
	file := &StorageFile{}
	err := GetDB().Model(&StorageFile{}).Where("id=? and user_id=?", id, swu.Owner.ID).First(file).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &utils.ErrResourceNotFound
		}
		return &utils.ErrInternalServerError
	}
	err = GetDB().Where("id=? ", id).Or("parent_id = ?", id).Delete(StorageFile{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &utils.ErrResourceNotFound
		}
		return &utils.ErrInternalServerError
	}
	return nil
}
