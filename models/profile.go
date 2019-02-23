package models

import (
	"net/http"

	"github.com/Dudobird/dudo-server/utils"
	"github.com/jinzhu/gorm"
)

// Profile for user information
type Profile struct {
	gorm.Model
	UserID       uint `json:"user_id"`
	User         User
	Name         string `json:"name" `
	Phone        string `json:"phone"`
	MobilePhone  string `json:"mobil_phone"`
	Department   string `json:"department"`
	ProfileImage string `json:"profile_image"`
}

// GetUserProfile return user profile struct
func GetUserProfile(accountID uint) *Profile {
	profile := &Profile{}
	err := GetDB().Table("profiles").Where("user_id = ?", accountID).First(profile).Error
	if err != nil {
		return nil
	}
	return profile
}

// Validate will check the field of Profile and
// return true if everything is fine
func (profile *Profile) Validate() (int, string) {
	if profile.UserID <= 0 {
		return http.StatusNotFound, "user id not found"
	}
	if profile.Name != "" && (len(profile.Name) <= 3 || len(profile.Name) >= 20) {
		return http.StatusBadRequest, "name lenth must greate than 3 and less than 20"
	}
	// More validate need
	return http.StatusOK, ""
}

// Create save all profile data and send 201 message back to user
func (profile *Profile) Create() *utils.Message {
	if status, message := profile.Validate(); status != http.StatusOK {
		return utils.NewMessage(status, message)
	}
	GetDB().Create(profile)
	message := utils.NewMessage(http.StatusCreated, "profile create success")
	message.Data = ""
	return message
}
