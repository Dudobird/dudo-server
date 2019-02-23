package controllers

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/Dudobird/dudo-server/auth"
	"github.com/Dudobird/dudo-server/models"
	"github.com/Dudobird/dudo-server/utils"
)

// GetTopLevelFiles list all top level files when user login success
func GetTopLevelFiles(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.TokenContextKey).(uint)
	user := models.GetUser(userID)

	if user == nil {
		utils.JSONRespnseWithTextMessage(w, http.StatusNotFound, "user not found")
		return
	}

	swu := models.StorageFilesWithUser{
		Owner: user,
	}
	err := swu.GetTopFiles()
	if err != nil {
		log.Errorln(err)
		utils.JSONRespnseWithTextMessage(w, http.StatusServiceUnavailable, "request service not avaliable now")
		return
	}
	utils.JSONMessageWithData(w, 200, "", swu.Files)
	return
}
