package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Dudobird/dudo-server/models"
	"github.com/Dudobird/dudo-server/utils"
)

type authInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleID   int    `json:"roleid"`
}

// CreateUser will create a new user based received json object
func CreateUser(w http.ResponseWriter, r *http.Request) {
	data := authInfo{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		utils.JSONRespnseWithErr(w, &utils.ErrPostDataNotCorrect)
		return
	}
	message := models.SignUp(data.Email, data.Password, data.RoleID)
	utils.JSONResonseWithMessage(w, message)
}

// Login will get user email and password from json object
// if user authentication information is correct, send back 200
// else send 403 forbidden
func Login(w http.ResponseWriter, r *http.Request) {
	data := authInfo{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		utils.JSONRespnseWithErr(w, &utils.ErrPostDataNotCorrect)
		return
	}
	message := models.Login(data.Email, data.Password)
	utils.JSONResonseWithMessage(w, message)
}

// Logout will logout user and delete the token infomation
// if user token not correct, send 401 unauthorization
// else send 200 logout success
func Logout(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(utils.TokenContextKey).(string)
	message := models.Logout(user)
	utils.JSONResonseWithMessage(w, message)
}

type receivePasswordInfo struct {
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
}

// UpdatePassword will update user password
// and user must send the new password and old password together
// if user token not correct, send 401 unauthorization
// else if old password not correct send 403 forbidden
// else if old password is correct but new password validate fail it will send 400
// else send 200 update success
func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.TokenContextKey).(string)
	tempAccout := &receivePasswordInfo{}
	err := json.NewDecoder(r.Body).Decode(tempAccout)
	if err != nil {
		utils.JSONRespnseWithErr(w, &utils.ErrPostDataNotCorrect)
		return
	}
	message := models.UpdatePassword(userID, tempAccout.Password, tempAccout.NewPassword)
	utils.JSONResonseWithMessage(w, message)
}
