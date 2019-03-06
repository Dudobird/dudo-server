package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Dudobird/dudo-server/models"
	"github.com/Dudobird/dudo-server/utils"
)

// e2e test with real backend and storage
func TestPostFormDataToCreateFiles(t *testing.T) {
	app := GetTestApp()
	userResponse, _ := signUpTestUser(app)
	token := userResponse.Data.Token
	testCases := []struct {
		url          string
		formDataName string
		filePath     string
		folderID     string
		token        string
		statusCode   int
	}{
		{
			url:          "/api/upload/storage/notexist",
			formDataName: "uploadfile",
			filePath:     "./files/1.file",
			token:        token,
			statusCode:   404,
		},
		{
			url:          "/api/upload/storage/root",
			formDataName: "notcorrect",
			filePath:     "./files/1.file",
			token:        token,
			statusCode:   400,
		},
		{
			url:          "/api/upload/storage/root",
			formDataName: "uploadfile",
			filePath:     "./files/1.file",
			token:        token,
			statusCode:   201,
		},
	}
	for _, tc := range testCases {
		rr, _ := fileUploadRequest(tc.url, tc.formDataName, tc.filePath, tc.token)
		utils.Equals(t, tc.statusCode, rr.Code)
	}
	tearDownUser(app)
	tearDownStorages()
}
func TestCreateFolders(t *testing.T) {
	app := GetTestApp()
	userResponse, _ := signUpTestUser(app)
	token := userResponse.Data.Token
	setUpRealFiles(token)
	defer func() {
		tearDownUser(app)
		tearDownStorages()
	}()
	testCase := []struct {
		fileInfo   []byte
		statuscode int
		token      string
	}{
		{
			// put a new folder under animals
			fileInfo:   []byte(`{"is_dir":true,"file_name":"wildanimals","parent_id":"bfc5dd70-f4e5-4aed-aad9-a9da313c8076"}`),
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
		req, _ := http.NewRequest("POST", "/api/storages", bytes.NewBuffer(test.fileInfo))
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
	req, _ := http.NewRequest("GET", "/api/storages", nil)
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

func TestGetTopFiles(t *testing.T) {
	app := GetTestApp()
	userResponse, _ := signUpTestUser(app)
	token := userResponse.Data.Token
	setUpRealFiles(token)
	defer func() {
		tearDownUser(app)
		tearDownStorages()
	}()
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
	req.Header.Set("Authorization", "Bearer "+token)
	rr = httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	utils.Equals(t, 200, rr.Code)
	message = StoragesResponse{}
	if err := json.NewDecoder(rr.Body).Decode(&message); err != nil {
		utils.OK(t, err)
	}
	for _, data := range message.Data {
		utils.Equals(t, data.ParentID, "")
	}
	utils.Equals(t, 3, len(message.Data))
}

func TestListCurrentFileWithID(t *testing.T) {
	app := GetTestApp()
	userResponse, _ := signUpTestUser(app)
	token := userResponse.Data.Token
	folders, _ := setUpRealFiles(token)
	defer func() {
		tearDownUser(app)
		tearDownStorages()
	}()
	testCase := []struct {
		id         string
		statuscode int
		token      string
		name       string
	}{
		{
			id:         folders["files"].ID,
			statuscode: 200,
			token:      token,
			name:       "files",
		},
		{
			id:         folders["empty"].ID,
			statuscode: 200,
			token:      token,
			name:       "empty",
		},
		{
			id:         folders["backup"].ID,
			statuscode: 200,
			token:      token,
			name:       "backup",
		},
		{
			id:         "not-exist",
			statuscode: 404,
			token:      token,
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
}

func TestListChildFilesWithID(t *testing.T) {
	app := GetTestApp()
	userResponse, _ := signUpTestUser(app)
	token := userResponse.Data.Token
	folders, _ := setUpRealFiles(token)
	defer func() {
		tearDownUser(app)
		tearDownStorages()
	}()
	testCase := []struct {
		id         string
		statuscode int
		token      string
		length     int
	}{
		{
			id:         folders["files"].ID,
			statuscode: 200,
			token:      token,
			length:     3,
		},
		{
			id:         folders["empty"].ID,
			statuscode: 200,
			token:      token,
			length:     0,
		},
		{
			id:         folders["backup"].ID,
			statuscode: 200,
			token:      token,
			length:     1,
		},
		{
			id:         "not-exist",
			statuscode: 200,
			token:      token,
			length:     0,
		},
		{
			id:         "bfc5dd70-f4e5-4aed-aad9-a9da313c8076",
			statuscode: 401,
			token:      "not-correct-token",
			length:     0,
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
			utils.Equals(t, test.length, len(message.Data))
		}
	}
}

func TestDeleteFilesWithID(t *testing.T) {
	app := GetTestApp()
	userResponse, _ := signUpTestUser(app)
	token := userResponse.Data.Token
	folders, files := setUpRealFiles(token)
	defer func() {
		tearDownUser(app)
		tearDownStorages()
	}()
	var counter int
	models.GetDB().Model(models.StorageFile{}).Count(&counter)
	utils.Equals(t, 7, counter)
	testCase := []struct {
		id         string
		name       string
		statuscode int
		token      string
	}{
		{
			// delete empty will remove this empty folder
			name:       "empty",
			id:         folders["empty"].ID,
			statuscode: 200,
			token:      token,
		},
		{
			// delete empty will remove this empty folder
			name:       "backup",
			id:         folders["backup"].ID,
			statuscode: 200,
			token:      token,
		},
		{
			// delete empty will remove this empty folder
			id:         files["1.file"].ID,
			name:       "1.file",
			statuscode: 200,
			token:      token,
		},
		{
			// delete empty will remove this empty folder
			id:         files["2.file"].ID,
			name:       "2.file",
			statuscode: 200,
			token:      token,
		},
		{
			id:         "not-exist",
			statuscode: 200,
			token:      token,
		},
		{
			id:         folders["backup"].ID,
			statuscode: 401,
			token:      "not-correct-token",
		},
		{
			id:         folders["backup"].ID,
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
			message := StoragesDeleteResponse{}
			if err := json.NewDecoder(rr.Body).Decode(&message); err != nil {
				utils.OK(t, err)
			}
		}
	}

	// check storage file no longer exist
	checkTestCases := []struct {
		id         string
		statuscode int
		name       string
		token      string
	}{
		{
			// check if delete success
			id:         folders["empty"].ID,
			statuscode: 404,
			name:       "empty",
			token:      token,
		},
		{
			id:         folders["backup"].ID,
			statuscode: 404,
			name:       "backup",
			token:      token,
		},
		{
			id:         files["1.file"].ID,
			statuscode: 404,
			name:       "1.file",
			token:      token,
		},
		{
			id:         files["2.file"].ID,
			statuscode: 404,
			token:      token,
		},
		{
			id:         files["3.file"].ID,
			statuscode: 200,
			name:       "3.file",
			token:      token,
		},
	}
	for _, test := range checkTestCases {
		req, _ := http.NewRequest("GET", "/api/storage/"+test.id, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+test.token)
		rr := httptest.NewRecorder()
		app.Router.ServeHTTP(rr, req)
		utils.Equals(t, test.statuscode, rr.Code)
	}
	models.GetDB().Model(models.StorageFile{}).Count(&counter)
	utils.Equals(t, 2, counter)
}

func TestDownloadFilesWithID(t *testing.T) {
	app := GetTestApp()
	userResponse, _ := signUpTestUser(app)
	token := userResponse.Data.Token
	folders, files := setUpRealFiles(token)
	defer func() {
		tearDownUser(app)
		tearDownStorages()
	}()
	testCases := []struct {
		fileName   string
		id         string
		statusCode int
		token      string
		savePath   string
		content    string
	}{
		{
			fileName:   "1.file",
			id:         files["1.file"].ID,
			statusCode: 200,
			token:      token,
			savePath:   "tmp-1.file",
			content:    "this is 1.file",
		},
		{
			fileName:   "2.file",
			id:         files["2.file"].ID,
			statusCode: 200,
			token:      token,
			savePath:   "tmp-2.file",
			content:    "this is 2.file",
		},
		{
			fileName:   "3.file",
			id:         files["3.file"].ID,
			statusCode: 200,
			token:      token,
			savePath:   "tmp-3.file",
			content:    "this is 3.file",
		},
		{
			fileName:   "2.file",
			id:         files["2.file"].ID,
			statusCode: 401,
			token:      "notexisttoken",
			savePath:   "this is 2.file",
		},
		{
			fileName:   "files",
			id:         folders["files"].ID,
			statusCode: 400,
			token:      token,
			savePath:   "tmp-fils.zip",
		},
	}
	for _, tc := range testCases {
		req, _ := http.NewRequest("GET", "/api/download/storage/"+tc.id, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tc.token)
		rr := httptest.NewRecorder()
		app.Router.ServeHTTP(rr, req)
		utils.Equals(t, tc.statusCode, rr.Code)
		if tc.statusCode == rr.Code {
			output, err := os.Create(tc.savePath)
			if err != nil {
				log.Printf("create temp file fail:%s", err)
				t.FailNow()
			}
			_, err = io.Copy(output, rr.Body)
			if err != nil {
				log.Printf("transfer file fail:%s", err)
				t.FailNow()
			}
			output.Close()
			if tc.content != "" {
				dat, _ := ioutil.ReadFile(tc.savePath)
				utils.Equals(t, tc.content, string(dat))
			}
			os.Remove(tc.savePath)
		}
	}

}
