package models

import (
	"log"

	"github.com/jinzhu/gorm"
)

// Storage for store user files
type Storage struct {
	gorm.Model
	UserID uint
	// filetype like exe jpeg
	RawStorageInfo
}

//RawStorageInfo after upload success ,each file will generate one raw info
type RawStorageInfo struct {
	FileExtention string `json:"file_extention"`
	FileName      string `json:"file_name"`
	FileLevel     uint   `json:"level"`
	Bucket        string `json:"bucket"`
	Path          string `json:"path"`
	IsDir         bool   `json:"isdir"`
}

// StoragesWithUser for controller
type StoragesWithUser struct {
	Owner    *User
	Storages []*Storage
}

// Save files from raw info
func (swu *StoragesWithUser) Save(files []RawStorageInfo) error {
	for _, file := range files {
		swu.Storages = append(swu.Storages, &Storage{RawStorageInfo: file, UserID: swu.Owner.ID})
	}
	err := GetDB().Model(swu.Owner).Association("Storages").Append(swu.Storages).Error
	return err
}

// GetTopFiles get user first level of files
func (swu *StoragesWithUser) GetTopFiles() error {
	storages := []*Storage{}
	err := GetDB().Model(swu.Owner).Where("file_level = ?", 0).Related(&storages).Error
	if err != nil {
		return err
	}
	log.Printf("Len of storage = %d", len(storages))
	swu.Storages = storages
	return nil
}
