package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// ShareFiles share file for unauthenticated use
type ShareFiles struct {
	gorm.Model
	FileID string    `json:"file_id"`
	Expire time.Time `json:"expire"`
	UserID string    `json:"user_id"`
}
