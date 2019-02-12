package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/zhangmingkai4315/dudo-server/models"
	"github.com/zhangmingkai4315/dudo-server/utils"
)

// CreateAccount will create a new user based received json object
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		utils.JSONRespnseWithTextMessage(w, http.StatusBadRequest, "request data invalid")
		return
	}
	message := account.Create()
	utils.JSONResonseWithMessage(w, message)
}

// Login will get user email and password from json object
// return user account information when success, or send back
// some error message
func Login(w http.ResponseWriter, r *http.Request) {
	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		utils.JSONRespnseWithTextMessage(w, http.StatusBadRequest, "request data invalid")
		return
	}
	message := models.Login(account.Email, account.Password)
	utils.JSONResonseWithMessage(w, message)
}

// Logout will logout user and delete the token infomation
func Logout(w http.ResponseWriter, r *http.Request) {

}

// Refresh token will refresh user token reset the expire date
func Refresh(w http.ResponseWriter, r *http.Request) {

}
