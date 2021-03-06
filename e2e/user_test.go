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

func TestCreateUser(t *testing.T) {
	app := GetTestApp()
	defer func() {
		tearDownUser(app)
	}()
	var users = []struct {
		post       []byte
		statuscode int
	}{
		{
			post:       []byte(`{"email":"test@example.com","password":"123456"}`),
			statuscode: http.StatusCreated,
		},
		{
			post:       []byte(`{"email":"test@example.com","password":""}`),
			statuscode: http.StatusBadRequest,
		},
		{
			post:       []byte(`{"email":"test@example.com"`),
			statuscode: http.StatusBadRequest,
		},
		{
			post:       []byte(`{"email":"test143432","password":"123456"}`),
			statuscode: http.StatusBadRequest,
		},
		{
			post:       []byte(``),
			statuscode: http.StatusBadRequest,
		},
		{
			post:       []byte(`{"password":"123456"}`),
			statuscode: http.StatusBadRequest,
		}, {
			post:       []byte(`{"email":"test@example.com","password":"2342253"}`),
			statuscode: http.StatusBadRequest,
		},
	}
	for _, user := range users {
		req, _ := http.NewRequest("POST", "/api/auth/signup", bytes.NewBuffer(user.post))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		app.Router.ServeHTTP(rr, req)
		utils.Equals(t, user.statuscode, rr.Code)
		if user.statuscode == http.StatusCreated {
			message := UserResponse{}
			if err := json.NewDecoder(rr.Body).Decode(&message); err != nil {
				utils.OK(t, err)
			}
			utils.Equals(t, message.Status, rr.Code)
			utils.Assert(t, message.Data.Password == "", "password not empty")
			utils.Assert(t, message.Data.Token != "", "token is empty")
		}
	}
}

func TestLoginUser(t *testing.T) {
	app := GetTestApp()
	_, err := signUpTestUser(app)
	defer func() {
		tearDownUser(app)
	}()
	utils.Equals(t, err, nil)
	var users = []struct {
		post       []byte
		statuscode int
	}{
		{
			post:       []byte(`{"email":"test@example.com","password":"123456"}`),
			statuscode: http.StatusOK,
		},
		{
			post:       []byte(`{"email":"test@example.com","password":"notcorrect"}`),
			statuscode: http.StatusForbidden,
		},
		{
			post:       []byte(`{"email":"test@example.com"`),
			statuscode: http.StatusBadRequest,
		},
		{
			post:       []byte(`{"email":"test143432","password":"123456"}`),
			statuscode: http.StatusBadRequest,
		},
		{
			post:       []byte(``),
			statuscode: http.StatusBadRequest,
		},
		{
			post:       []byte(`{"password":"123456"}`),
			statuscode: http.StatusBadRequest,
		}, {
			// with not correct password it will return forbidden
			post:       []byte(`{"email":"test@example.com","password":"2342253"}`),
			statuscode: http.StatusForbidden,
		},
	}
	for _, u := range users {
		req, _ := http.NewRequest("POST", "/api/auth/signin", bytes.NewBuffer(u.post))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		app.Router.ServeHTTP(rr, req)
		utils.Equals(t, u.statuscode, rr.Code)
		if u.statuscode == http.StatusOK {
			message := UserResponse{}
			if err := json.NewDecoder(rr.Body).Decode(&message); err != nil {
				utils.OK(t, err)
			}
			utils.Equals(t, message.Status, rr.Code)
			utils.Assert(t, message.Data.Password == "", "password not empty")
			utils.Assert(t, message.Data.Token != "", "token is empty")
		}
	}
}

func TestLogout(t *testing.T) {
	response, err := signUpTestUser(app)
	defer func() {
		tearDownUser(app)
	}()
	utils.Equals(t, err, nil)
	testtoken := response.Data.Token
	var testCases = []struct {
		token      string
		statuscode int
	}{
		{token: testtoken,
			statuscode: http.StatusOK,
		},
		{token: "",
			statuscode: http.StatusUnauthorized,
		},
	}

	for _, test := range testCases {
		req, _ := http.NewRequest("GET", "/api/auth/logout", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+test.token)
		rr := httptest.NewRecorder()
		app.Router.ServeHTTP(rr, req)
		utils.Equals(t, test.statuscode, rr.Code)
	}

}

func TestUpdatePassword(t *testing.T) {
	response, err := signUpTestUser(app)
	utils.Equals(t, err, nil)
	testtoken := response.Data.Token

	var testCases = []struct {
		postJSONString []byte
		statuscode     int
	}{

		{
			postJSONString: []byte(`{"new_password":"testpassword","password":"notcorrect"}`),
			statuscode:     http.StatusForbidden,
		},
		{
			postJSONString: []byte(`{"password":"123456"}`),
			statuscode:     http.StatusBadRequest,
		},
		{
			postJSONString: []byte(`{"new_password":"","password":"123456"}`),
			statuscode:     http.StatusBadRequest,
		},
		{
			postJSONString: []byte(`{"new_password":"testpassword","password":"123456"}`),
			statuscode:     http.StatusOK,
		},
	}

	for _, tc := range testCases {
		req, _ := http.NewRequest("POST", "/api/auth/password", bytes.NewBuffer(tc.postJSONString))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testtoken)
		rr := httptest.NewRecorder()
		app.Router.ServeHTTP(rr, req)
		utils.Equals(t, tc.statuscode, rr.Code)
		if tc.statuscode == http.StatusOK {
			newUser := &models.User{
				Email:    "test@example.com",
				Password: "testpassword",
			}
			user := newUser.ToJSONBytes()
			req, _ := http.NewRequest("POST", "/api/auth/signin", bytes.NewBuffer(user))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+testtoken)
			rr := httptest.NewRecorder()
			app.Router.ServeHTTP(rr, req)
			utils.Equals(t, http.StatusOK, rr.Code)
		}
	}
	tearDownUser(app)
}
