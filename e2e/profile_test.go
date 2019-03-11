package e2e

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Dudobird/dudo-server/models"

	"github.com/Dudobird/dudo-server/utils"
)

type ProfileResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    models.Profile
}

func TestGetUserProfile(t *testing.T) {
	app := GetTestApp()
	user, _ := signUpTestUser(app)
	defer func() {
		tearDownUser(app)
	}()
	token := user.Data.Token
	var testCases = []struct {
		token      string
		statuscode int
	}{
		{
			token:      token,
			statuscode: http.StatusOK,
		},
		{
			token:      "",
			statuscode: http.StatusUnauthorized,
		},
	}
	for _, tc := range testCases {
		req, _ := http.NewRequest("GET", "/api/profile", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tc.token)
		rr := httptest.NewRecorder()
		app.Router.ServeHTTP(rr, req)
		utils.Equals(t, tc.statuscode, rr.Code)
		if tc.statuscode == http.StatusOK {

			message := ProfileResponse{}
			if err := json.NewDecoder(rr.Body).Decode(&message); err != nil {
				utils.OK(t, err)
			}
			utils.Equals(t, message.Status, rr.Code)
			utils.Assert(t, message.Data.UserID == user.Data.ID, "user id is not same")
			utils.Assert(t, message.Data.UsageDiskSize == uint64(0), "init disk usage is not zero")
		}
	}
}

func TestUpdateUserProfile(t *testing.T) {
	app := GetTestApp()
	user, _ := signUpTestUser(app)
	defer func() {
		tearDownUser(app)
	}()
	token := user.Data.Token
	var testCases = []struct {
		token      string
		updateJSON []byte
		statuscode int
		expect     models.Profile
	}{
		{
			token:      token,
			updateJSON: []byte(`{"name":"Mike","disk_limit":0}`),
			statuscode: http.StatusOK,
		},
		{
			token:      "",
			updateJSON: []byte(`{"name":"Alice","disk_limit":0}`),
			statuscode: http.StatusUnauthorized,
		},
		{
			token:      token,
			updateJSON: []byte(`{"mobile_phone":"+86100000000","disk_limit":0}`),
			statuscode: http.StatusOK,
		},
		{
			token:      token,
			updateJSON: []byte(`{"department":"ops","disk_limit":0}`),
			statuscode: http.StatusOK,
		},

		{
			token:      token,
			updateJSON: []byte(`{"phone":"010-2342243","usage_disk_size":10}`),
			statuscode: http.StatusOK,
		},
	}
	for _, tc := range testCases {
		req, _ := http.NewRequest("PUT", "/api/profile", bytes.NewBuffer(tc.updateJSON))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tc.token)
		rr := httptest.NewRecorder()
		app.Router.ServeHTTP(rr, req)
		utils.Equals(t, tc.statuscode, rr.Code)
		if tc.statuscode == http.StatusOK {
			message := ProfileResponse{}
			if err := json.NewDecoder(rr.Body).Decode(&message); err != nil {
				utils.OK(t, err)
			}
			log.Printf("%+v", message.Data)
			utils.Equals(t, message.Status, rr.Code)
			utils.Assert(t, message.Data.UserID == user.Data.ID, "user id is not same")
			utils.Assert(t, message.Data.UsageDiskSize == uint64(0), "init disk usage is not zero")
		}
	}

	req, _ := http.NewRequest("GET", "/api/profile", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	utils.Equals(t, http.StatusOK, rr.Code)
	message := ProfileResponse{}
	if err := json.NewDecoder(rr.Body).Decode(&message); err != nil {
		utils.OK(t, err)
	}
	utils.Equals(t, message.Status, rr.Code)
	utils.Assert(t, message.Data.UserID == user.Data.ID, "user id is not same")
	utils.Assert(t, message.Data.UsageDiskSize == uint64(0), "user can change disksize?")
	utils.Assert(t, message.Data.Name == "Mike", "name not update")
	utils.Assert(t, message.Data.Phone == "010-2342243", "phone not update")
	utils.Assert(t, message.Data.MobilePhone == "+86100000000", "phone not update")
	utils.Assert(t, message.Data.Department == "ops", "department not update")
}
