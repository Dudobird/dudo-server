package models

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

var db *gorm.DB

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Panicf("Error: %s", err)
	}
	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dbPort := os.Getenv("db_port")
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, dbHost, dbPort, dbName)

	log.Printf("Info: connection uri %s\n", dbURI)
	conn, err := gorm.Open("mysql", dbURI)
	if err != nil {
		log.Panicf("Error: %s", err)
	}
	db = conn
	db.Debug().AutoMigrate(&Account{})
}

// GetDB will return a local db variable which init before
func GetDB() *gorm.DB {
	return db
}
