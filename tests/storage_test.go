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

func uploadFiles() {
	user := models.GetUserWithEmail(testUser.Email)
	if user == nil {
		log.Println("Error can't find test users")
	}
	files := models.StorageFilesWithUser{
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
	deleteFiles()
	deleteTestUser(app)
}

func TestUploadFilesAndListFilesWithID(t *testing.T) {
	response, err := signUpTestUser(app)
	utils.Equals(t, err, nil)
	testtoken := response.Data.Token
	uploadFiles()

	req, _ := http.NewRequest("GET", "/api/storages", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testtoken)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	message := StoragesResponse{}
	json.NewDecoder(rr.Body).Decode(&message)

	id := message.Data[0].ID

	testCase := []struct {
		id         string
		statuscode int
		token      string
	}{
		{
			id:         id,
			statuscode: 200,
			token:      testtoken,
		},
		{
			id:         "not-exist",
			statuscode: 400,
			token:      testtoken,
		},
		{
			id:         id,
			statuscode: 401,
			token:      "not-correct-token",
		},
		{
			id:         id,
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
			message := StoragesResponse{}
			if err := json.NewDecoder(rr.Body).Decode(&message); err != nil {
				utils.OK(t, err)
			}
			utils.Equals(t, 0, len(message.Data))
		}
	}
	deleteFiles()
	deleteTestUser(app)
}

func TestUploadFilesAndDeleteFilesWithID(t *testing.T) {
	response, err := signUpTestUser(app)
	utils.Equals(t, err, nil)
	testtoken := response.Data.Token
	uploadFiles()

	req, _ := http.NewRequest("GET", "/api/storages", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testtoken)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	message := StoragesResponse{}
	json.NewDecoder(rr.Body).Decode(&message)

	id := message.Data[0].ID

	testCase := []struct {
		id         string
		statuscode int
		token      string
	}{
		{
			id:         id,
			statuscode: 200,
			token:      testtoken,
		},
		{
			id:         "not-exist",
			statuscode: 400,
			token:      testtoken,
		},
		{
			id:         id,
			statuscode: 401,
			token:      "not-correct-token",
		},
		{
			id:         id,
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
			message := StoragesResponse{}
			if err := json.NewDecoder(rr.Body).Decode(&message); err != nil {
				utils.OK(t, err)
			}
			utils.Equals(t, 0, len(message.Data))
		}
	}
	deleteFiles()
	deleteTestUser(app)
}
