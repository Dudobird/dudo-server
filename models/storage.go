package models

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

// StorageFile for store user files
type StorageFile struct {
	ID        string `json:"id" gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	RawStorageFileInfo
	UserID   uint           `json:"user_id"`
	SubFiles []*StorageFile `gorm:"many2many:subfiles;association_jointable_foreignkey:subfile_id"`
}

// BeforeCreate save the uuid as file id
func (file *StorageFile) BeforeCreate(scope *gorm.Scope) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		log.Errorf("created uuid error: %v\n", err)
	}
	return scope.SetColumn("Id", uuid.String())
}

//RawStorageFileInfo after upload success ,each file will generate one raw info
type RawStorageFileInfo struct {
	FileExtention string `json:"file_extention"`
	FileName      string `json:"file_name" gorm:"not null;index:idx_file_name"`
	Bucket        string `json:"bucket" gorm:"not null;index:idx_file_bucket"`
	Path          string `json:"path" gorm:"not null"`
	IsDir         bool   `json:"isdir" gorm:"not null"`
	IsTopLevel    bool   `json:"is_top_level" gorm:"not null"`
}

// StorageFilesWithUser for controller
type StorageFilesWithUser struct {
	Owner *User
	Files []*StorageFile
}

// Save files from raw info
func (swu *StorageFilesWithUser) Save(files []RawStorageFileInfo) error {
	for _, file := range files {
		swu.Files = append(swu.Files, &StorageFile{RawStorageFileInfo: file})
	}
	err := GetDB().Model(swu.Owner).Association("Files").Append(swu.Files).Error
	return err
}

// GetTopFiles get user first level of files
func (swu *StorageFilesWithUser) GetTopFiles() error {
	files := []*StorageFile{}
	err := GetDB().Model(swu.Owner).Where("is_top_level = ?", true).Related(&files).Error
	if err != nil {
		return err
	}
	swu.Files = files
	return nil
}
