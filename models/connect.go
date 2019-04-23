package models

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Dudobird/dudo-server/config"
	"github.com/jinzhu/gorm"

	// mysql driver init
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

// InitDBConnection connect database
func InitDBConnection() (*gorm.DB, error) {
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
	log.Infof("try to connect database : %s:%s", dbHost, dbPort)
	db, err = gorm.Open("mysql", dbURI)
	if err != nil {
		log.Errorf("connect database error:%s", err)
		return nil, err
	}
	log.Infoln("connect database success")
	log.Infoln("start database automigrate")
	db.AutoMigrate(&User{}, &Profile{}, &StorageFile{}, &ShareFiles{}, &Role{})
	log.Infoln("database auto migrate success")

	// insert default data
	InsertDefaultRoles(db)

	return db, db.DB().Ping()
}

// InsertDefaultData insert or update default data
func InsertDefaultData(db *gorm.DB) error {
	// insert roles data
	err := InsertDefaultRoles(db)
	if err != nil {
		return err
	}
	return nil
}

// GetDB will return a local db variable which init before
func GetDB() *gorm.DB {
	if db == nil {
		db, err := InitDBConnection()
		if err != nil {
			panic(err)
		}
		return db
	}
	return db
}
