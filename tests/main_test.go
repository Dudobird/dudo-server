package tests

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	defer func() {
		cleanTables(app)
	}()
	app := GetTestApp()
	createTables(app)
	code := m.Run()
	app.DB.Close()
	os.Exit(code)
}
