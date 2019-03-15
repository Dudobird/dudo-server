package models

import (
	"time"
)

// ShareFiles share file for unauthenticated use
type ShareFiles struct {
	ID        string    `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at" gorm:"DEFAULT:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"DEFAULT:current_timestamp"`
	DeletedAt *time.Time

	StorageFile StorageFile `gorm:"foreignkey:FileID;auto_preload"`
	FileID      string      `json:"file_id"`
	Expire      time.Time   `json:"expire"`
	Description string      `json:"description" gorm:"not null;default:''"`
	UserID      string      `json:"user_id"`
}
