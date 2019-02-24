package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Dudobird/dudo-server/core"
	"github.com/Dudobird/dudo-server/models"
	"github.com/Dudobird/dudo-server/utils"
)

//                                    All Test Files
//             ___________________________|_________________________
//             |                             |                     |
//          animals                        test.zip               trees
//             |                                                    |
//    |——————————————————————————|                                  |
//    |                          |                                  |
//    cat.jpg                  dog.jpg                          pine.jpg
//

var rawFiles = []models.RawStorageFileInfo{
	{
		ID:         "4a447b2f-6947-478e-8207-20fd1f82d082",
		FileName:   "test.zip",
		Bucket:     "test",
		ParentID:   "",
		IsDir:      false,
		IsTopLevel: true,
	},
	{
		ID:         "bfc5dd70-f4e5-4aed-aad9-a9da313c8076",
		FileName:   "animals",
		Bucket:     "",
		ParentID:   "",
		IsDir:      true,
		IsTopLevel: true,
	},
	{
		ID:         "faeea8e1-3d9f-40c5-8097-121903d57339",
		FileName:   "trees",
		Bucket:     "",
		ParentID:   "",
		IsDir:      true,
		IsTopLevel: true,
	}, {
		ID:         "f2a5a7b9-e94c-4d0d-b48d-2597f41b199f",
		FileName:   "pine.jpg",
		Bucket:     "test",
		ParentID:   "faeea8e1-3d9f-40c5-8097-121903d57339",
		IsDir:      false,
		IsTopLevel: false,
	},
	{
		ID:         "6195f2f6-e12d-4bb7-a125-793a939caf6e",
		FileName:   "cat.jpg",
		Bucket:     "test",
		ParentID:   "bfc5dd70-f4e5-4aed-aad9-a9da313c8076",
		IsDir:      false,
		IsTopLevel: false,
	}, {
		ID:         "6ab7058e-5d90-453b-bfc8-96f936ddd815",
		FileName:   "dog.jpg",
		Bucket:     "test",
		ParentID:   "bfc5dd70-f4e5-4aed-aad9-a9da313c8076",
		IsDir:      false,
		IsTopLevel: false,
	},
}

type StoragesResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []struct {
		ID         string `json:"id"`
		FileName   string `json:"file_name"`
		Bucket     string `json:"bucket"`
		Path       string `json:"path"`
		IsTopLevel bool   `json:"is_top_level"`
		IsDir      bool   `json:"is_dir"`
		ParentID   string `json:"parent_id"`
	}
}

type SingleStoragesResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID         string `json:"id"`
		FileName   string `json:"file_name"`
		Bucket     string `json:"bucket"`
		Path       string `json:"path"`
		IsTopLevel bool   `json:"is_top_level"`
		IsDir      bool   `json:"is_dir"`
		ParentID   string `json:"parent_id"`
	}
}

func setUpUser(app *core.App) string {
	response, _ := signUpTestUser(app)
	return response.Data.Token
}

func setUpFiles() {
	user := models.GetUserWithEmail(testUser.Email)
	if user == nil {
		log.Panicln("Error can't find test users")
	}
	files := models.StorageFilesWithUser{
		Owner: user,
	}
	if err := files.SaveFromRawFiles(rawFiles); err != nil {
		log.Panicln("Save test files fail:" + err.Error())
	}
}

func tearDownUser(app *core.App) {
	deleteTestUser(app)
}
func tearDownStorages() {
	models.GetDB().Unscoped().Model(&models.StorageFile{}).Delete(&models.StorageFile{})
}
func TestEmptyFiles(t *testing.T) {
	app := GetTestApp()
	testtoken := setUpUser(app)
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

	tearDownUser(app)
}

func TestGetTopFiles(t *testing.T) {
	app := GetTestApp()
	testtoken := setUpUser(app)
	setUpFiles()

	// test with no authentication infomation
	req, _ := http.NewRequest("GET", "/api/storages", nil)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	utils.Equals(t, 401, rr.Code)
	message := StoragesResponse{}
	if err := json.NewDecoder(rr.Body).Decode(&message); err != nil {
		utils.OK(t, err)
	}

	// test with authentication infomation
	req, _ = http.NewRequest("GET", "/api/storages", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testtoken)
	rr = httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	utils.Equals(t, 200, rr.Code)
	message = StoragesResponse{}
	if err := json.NewDecoder(rr.Body).Decode(&message); err != nil {
		utils.OK(t, err)
	}
	for _, data := range message.Data {
		utils.Equals(t, data.IsTopLevel, true)
		utils.Equals(t, data.ParentID, "")
	}
	utils.Equals(t, 3, len(message.Data))
	tearDownUser(app)
	tearDownStorages()
}

func TestListCurrentFileWithID(t *testing.T) {
	app := GetTestApp()
	testtoken := setUpUser(app)
	setUpFiles()

	testCase := []struct {
		id         string
		statuscode int
		token      string
		name       string
	}{
		{
			id:         "bfc5dd70-f4e5-4aed-aad9-a9da313c8076",
			statuscode: 200,
			token:      testtoken,
			name:       "animals",
		},
		{
			id:         "f2a5a7b9-e94c-4d0d-b48d-2597f41b199f",
			statuscode: 200,
			token:      testtoken,
			name:       "pine.jpg",
		},
		{
			id:         "not-exist",
			statuscode: 400,
			token:      testtoken,
		},
		{
			id:         "bfc5dd70-f4e5-4aed-aad9-a9da313c8076",
			statuscode: 401,
			token:      "not-correct-token",
		},
		{
			id:         "bfc5dd70-f4e5-4aed-aad9-a9da313c8076",
			statuscode: 401,
			token:      "",
		},
	}

	for _, test := range testCase {
		req, _ := http.NewRequest("GET", "/api/storage/"+test.id, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+test.token)
		rr := httptest.NewRecorder()
		app.Router.ServeHTTP(rr, req)
		utils.Equals(t, test.statuscode, rr.Code)
		if rr.Code == 200 {
			message := SingleStoragesResponse{}
			if err := json.NewDecoder(rr.Body).Decode(&message); err != nil {
				utils.OK(t, err)
			}
			utils.Equals(t, test.name, message.Data.FileName)
		}
	}
	tearDownUser(app)
	tearDownStorages()
}

func TestListChildFilesWithID(t *testing.T) {
	app := GetTestApp()
	testtoken := setUpUser(app)
	setUpFiles()

	testCase := []struct {
		id         string
		statuscode int
		token      string
	}{
		{
			id:         "bfc5dd70-f4e5-4aed-aad9-a9da313c8076",
			statuscode: 200,
			token:      testtoken,
		},
		{
			id:         "not-exist",
			statuscode: 400,
			token:      testtoken,
		},
		{
			id:         "bfc5dd70-f4e5-4aed-aad9-a9da313c8076",
			statuscode: 401,
			token:      "not-correct-token",
		},
		{
			id:         "bfc5dd70-f4e5-4aed-aad9-a9da313c8076",
			statuscode: 401,
			token:      "",
		},
	}

	for _, test := range testCase {
		req, _ := http.NewRequest("GET", "/api/storage/"+test.id+"/subfiles", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+test.token)
		rr := httptest.NewRecorder()
		app.Router.ServeHTTP(rr, req)
		utils.Equals(t, test.statuscode, rr.Code)
		if rr.Code == 200 {
			message := StoragesResponse{}
			if err := json.NewDecoder(rr.Body).Decode(&message); err != nil {
				utils.OK(t, err)
			}
			utils.Equals(t, 2, len(message.Data))
		}
	}
	tearDownUser(app)
	tearDownStorages()
}

func TestDeleteFilesWithID(t *testing.T) {
	app := GetTestApp()
	testtoken := setUpUser(app)
	setUpFiles()
	testCase := []struct {
		id         string
		statuscode int
		token      string
	}{
		{
			// delete animal folder will also delete all subfiles
			id:         "bfc5dd70-f4e5-4aed-aad9-a9da313c8076",
			statuscode: 200,
			token:      testtoken,
		},
		{
			id:         "not-exist",
			statuscode: 400,
			token:      testtoken,
		},
		{
			id:         "bfc5dd70-f4e5-4aed-aad9-a9da313c8076",
			statuscode: 401,
			token:      "not-correct-token",
		},
		{
			id:         "bfc5dd70-f4e5-4aed-aad9-a9da313c8076",
			statuscode: 401,
			token:      "",
		},
	}

	for _, test := range testCase {
		req, _ := http.NewRequest("DELETE", "/api/storage/"+test.id, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+test.token)
		rr := httptest.NewRecorder()
		app.Router.ServeHTTP(rr, req)
		utils.Equals(t, test.statuscode, rr.Code)
		if rr.Code == 200 {
			message := StoragesResponse{}
			if err := json.NewDecoder(rr.Body).Decode(&message); err != nil {
				utils.OK(t, err)
			}
		}
	}

	// check storage file no longer exist

	testCase = []struct {
		id         string
		statuscode int
		token      string
	}{
		{
			// delete animal folder will also delete all subfiles
			id:         "bfc5dd70-f4e5-4aed-aad9-a9da313c8076",
			statuscode: 404,
			token:      testtoken,
		},
		{
			id:         "6195f2f6-e12d-4bb7-a125-793a939caf6e",
			statuscode: 404,
			token:      testtoken,
		},
		{
			id:         "6ab7058e-5d90-453b-bfc8-96f936ddd815",
			statuscode: 404,
			token:      testtoken,
		},
	}
	for _, test := range testCase {
		req, _ := http.NewRequest("GET", "/api/storage/"+test.id, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+test.token)
		rr := httptest.NewRecorder()
		app.Router.ServeHTTP(rr, req)
		utils.Equals(t, test.statuscode, rr.Code)
	}
	tearDownUser(app)
	tearDownStorages()
}

func TestCreateFolders(t *testing.T) {
	app := GetTestApp()
	testtoken := setUpUser(app)
	setUpFiles()
	testCase := []struct {
		fileInfo   []byte
		statuscode int
		token      string
	}{
		{
			// put a new folder under animals
			fileInfo:   []byte(`{"is_dir":true,"is_top_level":false,"file_name":"wildanimals","parent_id":"bfc5dd70-f4e5-4aed-aad9-a9da313c8076"}`),
			statuscode: 201,
			token:      testtoken,
		},
		{
			// create a new folder in root
			fileInfo:   []byte(`{"is_dir":true,"is_top_level":true,"file_name":"people"}`),
			statuscode: 201,
			token:      testtoken,
		},

		{
			// top folder with parent_id not allow
			fileInfo:   []byte(`{"is_dir":true,"is_top_level":true,"file_name":"test_animals","parent_id":"bfc5dd70-f4e5-4aed-aad9-a9da313c8076"}`),
			statuscode: 400,
			token:      testtoken,
		},
		// {
		// 	// put a new folder under animals must not exist (random result? maybe need some way to flush or sync)
		// 	fileInfo:   []byte(`{"is_dir":true,"is_top_level":false,"file_name":"wildanimals","parent_id":"bfc5dd70-f4e5-4aed-aad9-a9da313c8076"}`),
		// 	statuscode: 400,
		// 	token:      testtoken,
		// },
	}

	for _, test := range testCase {
		req, _ := http.NewRequest("POST", "/api/storages", bytes.NewBuffer(test.fileInfo))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+test.token)
		rr := httptest.NewRecorder()
		app.Router.ServeHTTP(rr, req)
		utils.Equals(t, test.statuscode, rr.Code)
	}
	tearDownUser(app)
	tearDownStorages()
}
