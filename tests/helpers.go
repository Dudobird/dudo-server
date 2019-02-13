package tests

import (
	"log"

	"github.com/zhangmingkai4315/dudo-server/core"
	"github.com/zhangmingkai4315/dudo-server/models"
)

var app *core.App

func GetTestApp() *core.App {
	if app == nil {
		log.Println("app is nil")
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
