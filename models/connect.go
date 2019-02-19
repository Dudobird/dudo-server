package models

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
	// mysql driver
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/Dudobird/dudo-server/config"
)

var db *gorm.DB

// InitConnection connect database
func InitConnection() (*gorm.DB, error) {
	var err error
	config := config.GetConfig()
	username := config.Database.Username
	password := config.Database.Password
	dbName := config.Database.DBName
	dbHost := config.Database.Host
	dbPort := config.Database.Port
	dbURI := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username,
		password,
		dbHost,
		dbPort,
		dbName,
	)
	db, err = gorm.Open("mysql", dbURI)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	log.Infoln("connect database success")
	db.AutoMigrate(&Account{})
	return db, nil
}

// GetDB will return a local db variable which init before
func GetDB() *gorm.DB {
	return db
}
