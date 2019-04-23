package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Dudobird/dudo-server/models"
	"github.com/Dudobird/dudo-server/utils"
)

type AdminUsersResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []models.User
}

func TestAdminGetUsers(t *testing.T) {
	app := GetTestApp()
	admin, _ := signUpAdminUser(app)
	normalUser, _ := signUpTestUser(app)
	normalUserToken := normalUser.Data.Token
	defer func() {
		tearDownUser(app)
	}()
	token := admin.Data.Token
	var testCases = []struct {
		token       string
		queryString string
		statuscode  int
	}{
		{
			token:       token,
			queryString: "?page=0&size=10",
			statuscode:  http.StatusOK,
		},
		{
			token:       normalUserToken,
			queryString: "?page=0&size=10",
			statuscode:  http.StatusUnauthorized,
		},
		{
			token:       "",
			queryString: "",
			statuscode:  http.StatusUnauthorized,
		},
		{
			token:       "",
			queryString: "?page=0&size=10",
			statuscode:  http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		req, _ := http.NewRequest("GET", "/api/admin/users"+tc.queryString, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tc.token)
		rr := httptest.NewRecorder()
		app.Router.ServeHTTP(rr, req)
		utils.Equals(t, tc.statuscode, rr.Code)
		if tc.statuscode == http.StatusOK {
			message := AdminUsersResponse{}
			if err := json.NewDecoder(rr.Body).Decode(&message); err != nil {
				utils.OK(t, err)
			}
			utils.Equals(t, message.Status, rr.Code)
			utils.Assert(t, len(message.Data) > 0, "at least have one admin user ")
		}
	}
}
