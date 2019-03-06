package tests

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	app := GetTestApp()
	cleanTables(app)
	createTables(app)
	code := m.Run()
	app.DB.Close()
	os.Exit(code)
}
