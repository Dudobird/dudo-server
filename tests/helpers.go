package tests

import (
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
	&models.Account{},
	&models.Profile{},
}

// CreateTables create table automatic
func CreateTables(app *core.App) {
	app.DB.AutoMigrate(appModels...)
}

// CleanTables will drop all models tables
func CleanTables(app *core.App) {
	app.DB.DropTable(appModels...)
}
