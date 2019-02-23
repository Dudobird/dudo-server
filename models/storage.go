package models

import (
	"errors"
	"time"
)

// StorageFile for store user files
type StorageFile struct {
	CreatedAt time.Time `gorm:"DEFAULT:current_timestamp"`
	UpdatedAt time.Time `gorm:"DEFAULT:current_timestamp"`
	DeletedAt *time.Time
	RawStorageFileInfo
	UserID   uint           `json:"user_id"`
	SubFiles []*StorageFile `gorm:"many2many:subfiles;association_jointable_foreignkey:subfile_id"`
}

//RawStorageFileInfo after upload success ,each file will generate one raw info
type RawStorageFileInfo struct {
	ID         string `json:"id" gorm:"primary_key"`
	FileName   string `json:"file_name" gorm:"not null;index:idx_file_name"`
	Bucket     string `json:"bucket"`
	ParentID   string `json:"parent_id" gorm:"not null"`
	IsDir      bool   `json:"isdir" gorm:"not null"`
	IsTopLevel bool   `json:"is_top_level" gorm:"not null"`
	// remote minio storage path
	Path string `json:"path" gorm:"not null"`
}

// StorageFilesWithUser for controller
type StorageFilesWithUser struct {
	Owner *User
	Files []*StorageFile
}

// ListChildren list the file and all subfiles
func (swu *StorageFilesWithUser) ListChildren(parentID string) ([]*StorageFile, error) {
	file := &StorageFile{}
	err := GetDB().Model(&StorageFile{}).Where("id=?", parentID).First(file).Error
	if err != nil {
		return nil, err
	}
	if file.UserID != swu.Owner.ID {
		return nil, errors.New("Forbidden")
	}
	GetDB().Preload("SubFiles").First(file)
	return file.SubFiles, nil
}

// Save files from raw info
func (swu *StorageFilesWithUser) Save(files []RawStorageFileInfo) error {
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
func (swu *StorageFilesWithUser) GetTopFiles() error {
	files := []*StorageFile{}
	err := GetDB().Model(swu.Owner).Where("is_top_level = ?", true).Related(&files).Error
	if err != nil {
		return err
	}
	swu.Files = files
	return nil
}

// DeleteFileFrom delete file based on its id and delete all subfiles
func (swu *StorageFilesWithUser) DeleteFileFrom(id string) error {
	file := &StorageFile{}
	err := GetDB().Model(&StorageFile{}).Where("id=?", id).First(file).Error
	if err != nil {
		return err
	}
	if file.UserID != swu.Owner.ID {
		return errors.New("Forbidden")
	}
	if err = GetDB().Model(file).Association("SubFiles").Clear().Error; err != nil {
		return err
	}

	if err = GetDB().Where("id=?", id).Delete(StorageFile{}).Error; err != nil {
		return err
	}
	return nil
}
