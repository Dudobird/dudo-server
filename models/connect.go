package models

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/zhangmingkai4315/dudo-server/config"
)

var db *gorm.DB

func initConnection() {
	config := config.GetConfig()
	username := config.Database.Username
	password := config.Database.Password
	dbName := config.Database.DBName
	dbHost := config.Database.Host
	dbPort := config.Database.Port
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, dbHost, dbPort, dbName)
	conn, err := gorm.Open("mysql", dbURI)
	if err != nil {
		log.Fatalln(err)
	}
	log.Infoln("connect database success")
	db = conn
	db.Debug().AutoMigrate(&Account{})
}

// GetDB will return a local db variable which init before
func GetDB() *gorm.DB {
	if db == nil {
		initConnection()
	}
	return db
}
