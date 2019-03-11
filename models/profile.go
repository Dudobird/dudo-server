package models

import (
	"github.com/Dudobird/dudo-server/utils"
	"github.com/jinzhu/gorm"
)

// Profile for user information
type Profile struct {
	gorm.Model
	UserID       string `json:"user_id"`
	Name         string `json:"name" gorm:"unique_index:idx_profile_name"`
	Phone        string `json:"phone"`
	MobilePhone  string `json:"mobile_phone"`
	Department   string `json:"department"`
	ProfileImage string `json:"profile_image"`

	DiskLimit     uint64 `json:"disk_limit"`
	UsageDiskSize uint64 `json:"usage_disk_size"`
}

// GetUserProfile return user profile struct
func GetUserProfile(accountID uint) (*Profile, *utils.CustomError) {
	profile := &Profile{}
	err := GetDB().Table("profiles").Where("user_id = ?", accountID).First(profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &utils.ErrResourceNotFound
		}
		return nil, &utils.ErrInternalServerError
	}
	return profile, nil
}

// Validate will check the field of Profile and
// return true if everything is fine
func (profile *Profile) Validate() *utils.CustomError {
	if profile.UserID == "" {
		return &utils.ErrUserNotFound
	}
	if profile.Name != "" && (len(profile.Name) <= 3 || len(profile.Name) >= 20) {
		return &utils.ErrValidationForProfileName
	}
	// More validate need
	return nil
}

// Create save all profile data and send 201 message back to user
func (profile *Profile) Create() *utils.CustomError {
	if errWithCode := profile.Validate(); errWithCode != nil {
		return errWithCode
	}
	err := GetDB().Create(profile).Error
	if err != nil {
		return &utils.ErrInternalServerError
	}
	return &utils.NoCreateErr
}
