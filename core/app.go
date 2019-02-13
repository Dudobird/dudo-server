package core

import (
	"errors"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/zhangmingkai4315/dudo-server/models"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/zhangmingkai4315/dudo-server/config"
	"github.com/zhangmingkai4315/dudo-server/routers"
)

// App is the manager of this application
type App struct {
	configfile string
	Config     *config.Config
	Router     *mux.Router
	DB         *gorm.DB
}

// NewApp create a new App struct from config file
func NewApp(file string) *App {
	app := &App{}
	err := app.init(file)
	if err != nil {
		log.Fatal(err)
	}
	return app
}

// Run will start the serve
func (app *App) Run() {
	hostAndPort := app.Config.Application.ListenAt
	log.Println("server start listen at:", hostAndPort)
	err := http.ListenAndServe(hostAndPort, app.Router)
	if err != nil {
		log.Fatal(err)
	}
}

// Init load the config file and init the database connection
func (app *App) init(configFile string) (err error) {
	if configFile == "" {
		return errors.New("config file path is empty")
	}
	app.configfile = configFile
	config, err := config.LoadConfig(configFile)
	if err != nil {
		return
	}
	app.Config = config

	router, err := routers.LoadRouters()
	if err != nil {
		return
	}
	app.Router = router
	db, err := models.InitConnection()
	if err != nil {
		return
	}
	app.DB = db
	return
}
