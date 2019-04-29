package models

import (
	"encoding/json"

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

// MarshalJSON for transfer user to readable json
func (profile *Profile) MarshalJSON() ([]byte, error) {
	type AliasStruct Profile
	return json.Marshal(&struct {
		ReadableDiskLimit string `json:"readable_disk_limit"`
		ReadableDiskUsage string `json:"readable_disk_usage"`
		*AliasStruct
	}{
		ReadableDiskLimit: utils.GetReadableFileSize(float64(profile.DiskLimit)),
		ReadableDiskUsage: utils.GetReadableFileSize(float64(profile.UsageDiskSize)),
		AliasStruct:       (*AliasStruct)(profile),
	})
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

// ChangeUserStorageSize change user disk size
func ChangeUserStorageSize(id, readableSize string) error {
	profile := Profile{}
	if id == "" || readableSize == "" {
		return utils.ErrPostDataNotCorrect
	}
	err := GetDB().Table("profiles").Where("user_id = ?", id).First(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &utils.ErrResourceNotFound
		}
		return &utils.ErrInternalServerError
	}
	profile.DiskLimit = utils.GetFileSizeFromReadable(readableSize)
	return db.Unscoped().Where("user_id = ?", id).Save(&profile).Error

}

// GetUserProfile return user profile struct
func GetUserProfile(accountID string) (*Profile, *utils.CustomError) {
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
