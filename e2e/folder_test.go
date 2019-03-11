package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Dudobird/dudo-server/models"
	"github.com/Dudobird/dudo-server/utils"
)

func TestCreateFolders(t *testing.T) {
	app := GetTestApp()
	userResponse, _ := signUpTestUser(app)
	token := userResponse.Data.Token
	defer func() {
		tearDownUser(app)
	}()
	testCase := []struct {
		fileInfo   []byte
		statuscode int
		token      string
	}{
		{
			// put a new folder under animals
			fileInfo:   []byte(`{"is_dir":true,"file_name":"wildanimals","folder_id":"bfc5dd70-f4e5-4aed-aad9-a9da313c8076"}`),
			statuscode: 404,
			token:      token,
		},
		{
			// create a new folder in root
			fileInfo:   []byte(`{"is_dir":true,"file_name":"people"}`),
			statuscode: 201,
			token:      token,
		},
		{
			// folder name empty will reject
			fileInfo:   []byte(`{"is_dir":true,"file_name":""}`),
			statuscode: 400,
			token:      token,
		},
	}

	for _, test := range testCase {
		req, _ := http.NewRequest("POST", "/api/folders", bytes.NewBuffer(test.fileInfo))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+test.token)
		rr := httptest.NewRecorder()
		app.Router.ServeHTTP(rr, req)
		utils.Equals(t, test.statuscode, rr.Code)
		if rr.Code == 201 {
			type result struct {
				Result int
			}
			var r result
			// the user id must equal folder user id
			models.GetDB().Raw("select count(*) as result from (select id from (select id from users union all select user_id from storage_files) tb1 group by id) tb2;").Scan(&r)
			utils.Equals(t, r.Result, 1)
		}
	}

}

func TestEmptyFiles(t *testing.T) {
	app := GetTestApp()
	userResponse, _ := signUpTestUser(app)
	defer func() {
		tearDownUser(app)
	}()
	token := userResponse.Data.Token
	req, _ := http.NewRequest("GET", "/api/folders/root", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	utils.Equals(t, 200, rr.Code)
	message := StoragesResponse{}
	if err := json.NewDecoder(rr.Body).Decode(&message); err != nil {
		utils.OK(t, err)
	}
	utils.Equals(t, 0, len(message.Data))
}

func TestGetFolderFiles(t *testing.T) {
	app := GetTestApp()
	userResponse, _ := signUpTestUser(app)
	token := userResponse.Data.Token
	setUpRealFiles(token)
	defer func() {
		tearDownUser(app)
		tearDownStorages()
	}()
	testCase := []struct {
		folderPath string
		statuscode int
		token      string
	}{
		{
			folderPath: "root",
			statuscode: 401,
			token:      "",
		},
		{
			folderPath: "root",
			statuscode: 200,
			token:      token,
		},
	}

	for _, tc := range testCase {
		req, _ := http.NewRequest("GET", "/api/folders/"+tc.folderPath, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tc.token)
		rr := httptest.NewRecorder()
		app.Router.ServeHTTP(rr, req)
		utils.Equals(t, tc.statuscode, rr.Code)
		if rr.Code == 200 {
			message := StoragesResponse{}
			if err := json.NewDecoder(rr.Body).Decode(&message); err != nil {
				utils.OK(t, err)
			}
			for _, data := range message.Data {
				utils.Equals(t, data.FolderID, tc.folderPath)

			}
		}
	}
}

func TestRenameFileAndFolder(t *testing.T) {
	app := GetTestApp()
	userResponse, _ := signUpTestUser(app)
	token := userResponse.Data.Token
	_, files := setUpRealFiles(token)
	defer func() {
		tearDownUser(app)
		tearDownStorages()
	}()
	testCases := []struct {
		fileName   string
		id         string
		statusCode int
		token      string
		post       []byte
		newName    string
	}{
		{
			fileName:   "1.file",
			id:         files["1.file"].ID,
			statusCode: 200,
			token:      token,
			post:       []byte(`{"file_name":"changed_1.file"}`),
			newName:    "changed_1.file",
		},
		{
			fileName:   "1.file",
			id:         files["1.file"].ID,
			statusCode: 401,
			token:      "",
			post:       []byte(`{"file_name":"changed_1.file"}`),
			newName:    "changed_1.file",
		},
		{
			fileName:   "1.file",
			id:         files["1.file"].ID,
			statusCode: 200,
			token:      token,
			post:       []byte(`{"file_name":"1.file"}`),
			newName:    "1.file",
		},
		{
			fileName:   "1.file",
			id:         files["1.file"].ID,
			statusCode: 400,
			token:      token,
			post:       []byte(`{"file_name":"2.file"}`),
			newName:    "2.file",
		},
	}
	for _, tc := range testCases {
		req, _ := http.NewRequest("PUT", "/api/files/"+tc.id, bytes.NewBuffer(tc.post))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tc.token)
		rr := httptest.NewRecorder()
		app.Router.ServeHTTP(rr, req)
		utils.Equals(t, tc.statusCode, rr.Code)
		if rr.Code == 200 {
			message := SingleStoragesResponse{}
			if err := json.NewDecoder(rr.Body).Decode(&message); err != nil {
				utils.OK(t, err)
			}
			utils.Equals(t, message.Data.FileName, tc.newName)
		}
	}
}
