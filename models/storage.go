package models

import (
	"github.com/jinzhu/gorm"
)

// Storage for store user files
type Storage struct {
	gorm.Model
	UserID uint
	// filetype like exe jpeg
	FileExtention string `json:"file_extention"`
	FileName      string `json:"file_name"`
	FileLevel     uint   `json:"level"`
	Bucket        string `json:"bucket"`
	Path          string `json:"path"`
}
