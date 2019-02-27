package config

import (
	"testing"

	"github.com/Dudobird/dudo-server/utils"
)

var expectConfig = &Config{
	Application: application{
		ListenAt:     "127.0.0.1:8080",
		Token:        "thisisonlyfortest",
		TempFolder:   "temp",
		BucketPrefix: "dudotest",
	},
	Database: database{
		Type:     "mysql",
		Host:     "localhost",
		Port:     "3306",
		DBName:   "dudo",
		Password: "dudoadmin",
		Username: "dudouser",
	},
	Storage: storage{
		Server:    "localhost",
		Port:      "9000",
		AccessKey: "minio",
		SecretKey: "minio123",
	},
}

func TestLoadConfig(t *testing.T) {
	file := ""
	config, err := LoadConfig(file)
	utils.Assert(t, config == nil, "with empty file path LoadConfig() should return config==nil")
	utils.Assert(t, err != nil, "with empty file path LoadConfig() should return err!=nil")

	correctFile := "example.toml"
	config, err = LoadConfig(correctFile)
	utils.Assert(t, config != nil, "with correct file path LoadConfig() should return config==nil")
	utils.Assert(t, err == nil, "with correct file path LoadConfig() should return err!=nil")

	utils.Equals(t, expectConfig, config)
}

func TestGetConfig(t *testing.T) {
	correctFile := "example.toml"
	config, err := LoadConfig(correctFile)
	utils.Assert(t, config != nil, "with correct file path LoadConfig() should return config==nil")
	utils.Assert(t, err == nil, "with correct file path LoadConfig() should return err!=nil")

	utils.Equals(t, config, GetConfig())

	utils.Equals(t, expectConfig, config)

}
