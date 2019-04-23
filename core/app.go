package core

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Dudobird/dudo-server/models"
	"github.com/jinzhu/gorm"

	"github.com/Dudobird/dudo-server/config"
	"github.com/Dudobird/dudo-server/storage"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

// App is the manager of this application
type App struct {
	configfile     string
	Config         *config.Config
	Router         *mux.Router
	DB             *gorm.DB
	Storage        storage.Storage
	FullTempFolder string
}

var globalApp *App

// GetApp return global app
func GetApp() *App {
	return globalApp
}

// NewApp create a new App struct from config file
func NewApp(file string) *App {
	newApp := &App{}
	err := newApp.init(file)
	if err != nil {
		log.Fatal(err)
	}
	globalApp = newApp
	return newApp
}

// Run will start the serve
func (app *App) Run() {
	hostAndPort := app.Config.Application.ListenAt
	log.Println("server start listen at:", hostAndPort)
	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-FilePath"},
	})
	err := http.ListenAndServe(hostAndPort, c.Handler(app.Router))
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

	db, err := models.InitDBConnection()
	if err != nil {
		return
	}
	err = models.InsertDefaultData(db)
	if err != nil {
		return
	}
	app.DB = db
	app.Storage = storage.InitStorageManager()
	if err != nil {
		return
	}
	// Create temp file for file upload
	tempFolderName := config.Application.TempFolder
	fullTempPath, _ := filepath.Abs("." + string(filepath.Separator) + tempFolderName)
	if _, err = os.Stat(fullTempPath); os.IsNotExist(err) {
		err = os.MkdirAll(fullTempPath, 0755)
		if err != nil {
			return
		}
	}
	app.FullTempFolder = fullTempPath
	return
}
