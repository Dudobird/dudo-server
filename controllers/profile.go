package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/Dudobird/dudo-server/models"
	"github.com/Dudobird/dudo-server/utils"
	log "github.com/sirupsen/logrus"
)

// CreateProfile for create user profile
func CreateProfile(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(utils.TokenContextKey).(string)
	profile := &models.Profile{}

	err := json.NewDecoder(r.Body).Decode(profile)
	if err != nil {
		utils.JSONRespnseWithErr(w, &utils.ErrPostDataNotCorrect)
		return
	}
	profile.UserID = user
	utils.JSONRespnseWithErr(w, profile.Create())
	return
}

// GetProfile retrive user profile with user id
func GetProfile(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(utils.TokenContextKey).(string)
	profile := &models.Profile{}
	err := models.GetDB().Model(&models.Profile{}).Where("user_id = ?", user).First(profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.JSONRespnseWithErr(w, &utils.ErrUserNotFound)
			return
		}
		utils.JSONRespnseWithErr(w, &utils.ErrInternalServerError)
		return
	}
	utils.JSONMessageWithData(w, http.StatusOK, "", profile)
	return
}

// UpdateProfile retrive user profile with user id
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(utils.TokenContextKey).(string)

	profile := &models.Profile{}
	err := json.NewDecoder(r.Body).Decode(profile)
	currentUserProfile := &models.Profile{}
	err = models.GetDB().Model(&models.Profile{}).Where("user_id = ?", user).First(currentUserProfile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.JSONRespnseWithErr(w, &utils.ErrUserNotFound)
			return
		}
		log.Errorf("find user profile error: %s", err)
		utils.JSONRespnseWithErr(w, &utils.ErrInternalServerError)
		return
	}
	err = models.GetDB().Model(currentUserProfile).Omit("disk_limit", "usage_disk_size").Updates(profile).Error
	if err != nil {
		log.Errorf("update user profile error: %s", err)
		utils.JSONRespnseWithErr(w, &utils.ErrInternalServerError)
		return
	}
	models.GetDB().Model(&models.Profile{}).Save(currentUserProfile)
	utils.JSONMessageWithData(w, http.StatusOK, "", currentUserProfile)

	return
}
