package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/Dudobird/dudo-server/controllers"
	"github.com/Dudobird/dudo-server/core"
	"github.com/Dudobird/dudo-server/models"
	"github.com/Dudobird/dudo-server/routers"
	"github.com/Dudobird/dudo-server/utils"
	"github.com/gorilla/mux"
)

// test application
var app *core.App

// GetTestApp load test config return a new app
func GetTestApp() *core.App {
	if app == nil {
		app = core.NewApp("./config_test.toml")
		router, err := routers.LoadRouters()
		if err != nil {
			panic(err)
		}
		app.Router = router
	}
	return app
}

var appModels = []interface{}{
	&models.User{},
	&models.Profile{},
	&models.StorageFile{},
}

// createTables create table automatic
func createTables(app *core.App) {
	log.Println("create tables for test")
	app.DB.AutoMigrate(appModels...)
	log.Println("data migrate success")
}

// cleanTables will drop all models tables
func cleanTables(app *core.App) {
	app.DB.DropTable(appModels...)
	log.Println("clean all tables success")
}

func tearDownStorages() {
	models.GetDB().Unscoped().Model(&models.StorageFile{}).Delete(&models.StorageFile{})
	userID := strings.ToLower(strings.TrimLeft(UserID, "user_"))
	bucketName := fmt.Sprintf("dudotest-%s", userID)
	GetTestApp().Storage.CleanBucket(bucketName)
	log.Println("remove files data for database and storage success")
}

// var testUser models.User
var testUser = &models.User{
	Email:    "test@example.com",
	Password: "123456",
}

// UserID save created user id
var UserID string

// UserResponse save response from api response
type UserResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Email    string `json:"email"`
		Token    string `json:"token"`
		Password string `json:"password"`
		ID       string `json:"id"`
	}
}

func signUpTestUser(app *core.App) (*UserResponse, error) {
	user := testUser.ToJSONBytes()
	req, _ := http.NewRequest("POST", "/api/auth/signup", bytes.NewBuffer(user))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	message := UserResponse{}
	if http.StatusCreated == rr.Code {
		if err := json.NewDecoder(rr.Body).Decode(&message); err != nil {
			return nil, err
		}
		UserID = message.Data.ID
		return &message, nil
	}
	return nil, fmt.Errorf("sign up user fail with code %d", rr.Code)
}

func tearDownUser(app *core.App) {
	app.DB.Unscoped().Delete(&models.User{})
}

// StoragesResponse save the response infomation
// from server file list query api
type StoragesResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []struct {
		ID       string `json:"id"`
		FileName string `json:"file_name"`
		Bucket   string `json:"bucket"`
		Path     string `json:"path"`
		IsDir    bool   `json:"is_dir"`
		ParentID string `json:"parent_id"`
	}
}

// StoragesDeleteResponse save the response infomation
// from server file delete api
type StoragesDeleteResponse struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Data    []string `json:"data"`
}

// SingleStoragesResponse save the response infomation
// from server file query api for one file
type SingleStoragesResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID       string `json:"id"`
		FileName string `json:"file_name"`
		Bucket   string `json:"bucket"`
		Path     string `json:"path"`
		IsDir    bool   `json:"is_dir"`
		ParentID string `json:"parent_id"`
	}
}

//  save some files for test
//                        			 All  Test Files
//             ___________________________|_________________________
//             |                          |                         |
//           files                        empty                  backup
//             |                                                    |
//    |——————————————————————————|                                  |
//    |            |             |                                  |
//    1.file      2.file        3.file                            1.file

func setUpRealFiles(token string) (map[string]models.StorageFile, map[string]models.StorageFile) {
	// create folder
	testCase := []struct {
		fileInfo   []byte
		statuscode int
		token      string
	}{
		{
			// put a new folder under animals
			fileInfo:   []byte(`{"is_dir":true,"file_name":"files"}`),
			statuscode: 201,
			token:      token,
		},
		{
			// put a new folder under animals
			fileInfo:   []byte(`{"is_dir":true,"file_name":"empty"}`),
			statuscode: 201,
			token:      token,
		},
		{
			// put a new folder under animals
			fileInfo:   []byte(`{"is_dir":true,"file_name":"backup"}`),
			statuscode: 201,
			token:      token,
		},
	}

	for _, test := range testCase {
		req, _ := http.NewRequest("POST", "/api/storages", bytes.NewBuffer(test.fileInfo))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+test.token)
		rr := httptest.NewRecorder()
		app.Router.ServeHTTP(rr, req)
		if rr.Code != test.statuscode {
			log.Panicf("create folders fail: %s", rr.Body.String())
		}
	}
	folders := make(map[string]models.StorageFile)
	for _, name := range []string{"files", "empty", "backup"} {
		s := models.StorageFile{}
		err := models.GetDB().Model(models.StorageFile{}).Where("file_name = ?", name).First(&s).Error
		if err != nil {
			log.Panic(err)
		}
		folders[name] = s
	}
	realfiles := []struct {
		url          string
		formDataName string
		filePath     string
		token        string
		statusCode   int
	}{
		{
			url:          "/api/upload/storage/" + folders["files"].ID,
			formDataName: "uploadfile",
			filePath:     "./files/1.file",
			token:        token,
			statusCode:   201,
		},
		{
			url:          "/api/upload/storage/" + folders["backup"].ID,
			formDataName: "uploadfile",
			filePath:     "./files/1.file",
			token:        token,
			statusCode:   201,
		},
		{
			url:          "/api/upload/storage/" + folders["files"].ID,
			formDataName: "uploadfile",
			filePath:     "./files/2.file",
			token:        token,
			statusCode:   201,
		},
		{
			url:          "/api/upload/storage/" + folders["files"].ID,
			formDataName: "uploadfile",
			filePath:     "./files/3.file",
			token:        token,
			statusCode:   201,
		},
	}
	for _, tc := range realfiles {
		rr, err := fileUploadRequest(tc.url, tc.formDataName, tc.filePath, tc.token)
		if err != nil {
			log.Panic(err)
		}
		if rr.Code != http.StatusCreated {
			log.Panicf("statuscode = %d, body = %s", rr.Code, rr.Body.String())
		}
	}
	files := make(map[string]models.StorageFile)
	rawFiels := []models.StorageFile{}
	models.GetDB().Where("is_dir = ?", false).Find(&rawFiels)
	// only return three files in `files` folder
	for _, file := range rawFiels {

		if file.ParentID != folders["backup"].ID {

			files[file.FileName] = file
		}
	}
	return folders, files
}

func fileUploadRequest(url, paramName, path, token string) (*httptest.ResponseRecorder, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	file.Close()
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, fi.Name())
	if err != nil {
		return nil, err
	}
	part.Write(fileContents)
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	req := httptest.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()
	// handler := http.HandlerFunc(controllers.UploadFiles)
	ctx := req.Context()
	user := &models.User{}
	err = models.GetDB().Where("email = ?", testUser.Email).First(user).Error
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, utils.ContextToken("MyAppToken"), user.ID)
	req = req.WithContext(ctx)

	router := mux.NewRouter()
	router.HandleFunc("/api/upload/storage/{parentID}", controllers.UploadFiles)
	router.ServeHTTP(rr, req)
	return rr, nil
}
