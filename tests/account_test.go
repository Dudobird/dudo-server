package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEmptyAccountTable(t *testing.T) {
	app := GetTestApp()
	req, _ := http.NewRequest("GET", "/users", nil)
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	if http.StatusOK != rr.Code {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, rr.Code)
	}
}
