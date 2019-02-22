package tests

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Dudobird/dudo-server/models"

	"github.com/Dudobird/dudo-server/utils"
)

var rawFiles = []models.RawStorageInfo{
	{
		FileExtention: "jpg",
		FileName:      "cat",
		FileLevel:     0,
		Bucket:        "test",
		Path:          "cat.jpg",
	}, {
		FileExtention: "jpg",
		FileName:      "dog",
		FileLevel:     0,
		Bucket:        "test",
		Path:          "dog.jpg",
	},
}

type StoragesResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []struct {
		FileExtention string `json:"file_extention"`
		FileName      string `json:"file_name"`
		FileLevel     uint   `json:"level"`
		Bucket        string `json:"bucket"`
		Path          string `json:"path"`
	}
}

func uploadFiles() {
	user := models.GetUserWithEmail(testUser.Email)
	if user == nil {
		log.Println("Error can't find test users")
	}
	files := models.StoragesWithUser{
		Owner: user,
	}
	if err := files.Save(rawFiles); err != nil {
		log.Println("Save test files fail:" + err.Error())
	}
}

func deleteFiles() {
	user := models.GetUserWithEmail(testUser.Email)
	models.GetDB().Model(user).Association("Storages").Clear()
}
func TestNewUserWithEmptyStorage(t *testing.T) {
	response, err := signUpTestUser(app)
	utils.Equals(t, err, nil)
	testtoken := response.Data.Token
	req, _ := http.NewRequest("GET", "/api/storages", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testtoken)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	utils.Equals(t, 200, rr.Code)

	message := StoragesResponse{}
	if err := json.NewDecoder(rr.Body).Decode(&message); err != nil {
		utils.OK(t, err)
	}
	utils.Equals(t, 0, len(message.Data))
	deleteTestUser(app)
}

func TestUploadFilesAndGetFiles(t *testing.T) {
	response, err := signUpTestUser(app)
	utils.Equals(t, err, nil)
	testtoken := response.Data.Token
	uploadFiles()
	req, _ := http.NewRequest("GET", "/api/storages", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testtoken)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	utils.Equals(t, 200, rr.Code)

	message := StoragesResponse{}
	if err := json.NewDecoder(rr.Body).Decode(&message); err != nil {
		utils.OK(t, err)
	}
	utils.Equals(t, 2, len(message.Data))
	deleteFiles()
	deleteTestUser(app)
}
