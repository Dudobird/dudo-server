package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zhangmingkai4315/dudo-server/utils"
)

type AccountResponse struct {
	Status  int    `json:"status"`
	Message string `json:message`
	Data    struct {
		Email    string `json:"email"`
		Token    string `json:"token"`
		Password string `json:"password"`
	}
}

func TestCreateAccount(t *testing.T) {
	app := GetTestApp()
	var user = []byte(`{"email":"test@example.com","password":"123456"}`)
	req, _ := http.NewRequest("POST", "/api/auth/signup", bytes.NewBuffer(user))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	if http.StatusCreated != rr.Code {
		log.Println(rr.Body.String())
		utils.Equals(t, http.StatusCreated, rr.Code)
	}
	message := AccountResponse{}
	if err := json.NewDecoder(rr.Body).Decode(&message); err != nil {
		utils.OK(t, err)
	}
	utils.Equals(t, message.Status, rr.Code)
	utils.Assert(t, message.Data.Password == "", "password not empty")
	utils.Assert(t, message.Data.Token != "", "token is empty")
}
