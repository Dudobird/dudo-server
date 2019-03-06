package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/Dudobird/dudo-server/auth"
	"github.com/Dudobird/dudo-server/models"
	"github.com/Dudobird/dudo-server/utils"
)

// CreateProfile for create user profile
func CreateProfile(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(auth.TokenContextKey).(string)
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
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		utils.JSONRespnseWithErr(w, &utils.ErrPostDataNotCorrect)
		return
	}
	data, errWithCode := models.GetUserProfile(uint(id))
	if errWithCode != nil {
		utils.JSONRespnseWithErr(w, errWithCode)
		return
	}
	utils.JSONMessageWithData(w, http.StatusOK, "", data)
	return
}
