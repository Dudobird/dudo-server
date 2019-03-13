package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Dudobird/dudo-server/utils"
)

func TestCreateShareFile(t *testing.T) {
	app := GetTestApp()
	userResponse, _ := signUpTestUser(app)
	token := userResponse.Data.Token
	_, files := setUpRealFiles(token)
	defer func() {
		tearDownUser(app)
		tearDownStorages()
	}()
	testCase := []struct {
		postJSONString []byte
		statuscode     int
		token          string
	}{
		{
			postJSONString: []byte(fmt.Sprintf(`{"file_id":"%s","expire":7}`, files["1.file"].ID)),
			statuscode:     201,
			token:          token,
		},
		{
			postJSONString: []byte(fmt.Sprintf(`{"file_id":"%s","expire":7}`, files["1.file"].ID)),
			statuscode:     401,
			token:          "not-correct-token",
		},
		{
			postJSONString: []byte(`{"file_id":"not-exist","expire":7}`),
			statuscode:     404,
			token:          token,
		},
		{
			postJSONString: []byte(``),
			statuscode:     400,
			token:          token,
		},
	}

	for _, tc := range testCase {
		req, _ := http.NewRequest("POST", "/api/shares", bytes.NewBuffer(tc.postJSONString))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tc.token)
		rr := httptest.NewRecorder()
		app.Router.ServeHTTP(rr, req)
		utils.Equals(t, tc.statuscode, rr.Code)
		if rr.Code == 200 {
			message := struct {
				Data struct {
					Token string `json:"token"`
				}
			}{}
			if err := json.NewDecoder(rr.Body).Decode(&message); err != nil {
				utils.OK(t, err)
			}
			utils.Assert(t, len(message.Data.Token) > 0, "token must exist and not empty")
		}
	}
}
