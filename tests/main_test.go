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
	os.Exit(code)
}
