package config

import (
	"testing"

	"github.com/zhangmingkai4315/dudo-server/utils"
)

var expectConfig = &Config{
	Application: application{
		ListenAt: "127.0.0.1:8080",
		Token:    "thisisonlyfortest",
	},
	Database: database{
		Type:     "mysql",
		Host:     "localhost",
		Port:     "3306",
		DBName:   "dudo",
		Password: "dudoadmin",
		Username: "dudouser",
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
