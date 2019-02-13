package tests

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	app := GetTestApp()
	CreateTables(app)
	code := m.Run()
	CleanTables(app)
	app.DB.Close()
	os.Exit(code)
}
