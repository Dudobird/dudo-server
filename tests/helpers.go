package tests

import (
	"log"

	"github.com/Dudobird/dudo-server/core"
	"github.com/Dudobird/dudo-server/models"
)

var app *core.App

// GetTestApp load test config return a new app
func GetTestApp() *core.App {
	if app == nil {
		app = core.NewApp("./config_test.toml")
	}
	return app
}

var appModels = []interface{}{
	&models.User{},
	&models.Profile{},
}

// CreateTables create table automatic
func CreateTables(app *core.App) {
	log.Println("create tables for test")
	app.DB.AutoMigrate(appModels...)
}

// CleanTables will drop all models tables
func CleanTables(app *core.App) {
	log.Println("clean all tables")
	app.DB.DropTable(appModels...)
}
