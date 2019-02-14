package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zhangmingkai4315/dudo-server/models"

	"github.com/zhangmingkai4315/dudo-server/utils"
)

type AccountResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Email    string `json:"email"`
		Token    string `json:"token"`
		Password string `json:"password"`
	}
}

func TestCreateAccount(t *testing.T) {
	app := GetTestApp()
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
			// return status bad request becouse same email is already in use
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
			message := AccountResponse{}
			if err := json.NewDecoder(rr.Body).Decode(&message); err != nil {
				utils.OK(t, err)
			}
			utils.Equals(t, message.Status, rr.Code)
			utils.Assert(t, message.Data.Password == "", "password not empty")
			utils.Assert(t, message.Data.Token != "", "token is empty")
		}
	}

	app.DB.Unscoped().Delete(&models.Account{Email: "test@example.com"})
}

func TestLoginAccount(t *testing.T) {
	app := GetTestApp()
	user := []byte(`{"email":"test@example.com","password":"123456"}`)
	req, _ := http.NewRequest("POST", "/api/auth/signup", bytes.NewBuffer(user))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	utils.Equals(t, http.StatusCreated, rr.Code)

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
		req, _ = http.NewRequest("POST", "/api/auth/signin", bytes.NewBuffer(u.post))
		req.Header.Set("Content-Type", "application/json")
		rr = httptest.NewRecorder()
		app.Router.ServeHTTP(rr, req)
		utils.Equals(t, u.statuscode, rr.Code)
		if u.statuscode == http.StatusOK {
			message := AccountResponse{}
			if err := json.NewDecoder(rr.Body).Decode(&message); err != nil {
				utils.OK(t, err)
			}
			utils.Equals(t, message.Status, rr.Code)
			utils.Assert(t, message.Data.Password == "", "password not empty")
			utils.Assert(t, message.Data.Token != "", "token is empty")
		}
	}
	app.DB.Unscoped().Delete(&models.Account{Email: "test@example.com"})
}
