package config

import (
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

// Config store all config items
type Config struct {
	Database    database    `toml:"Database"`
	Application application `toml:"Application"`
}

type database struct {
	Type     string `toml:"type"`
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	DBName   string `toml:"dbname"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

type application struct {
	ListenAt string `toml:"listenAt"`
	Token    string `toml:"token"`
}

var config *Config

// LoadConfig return the config object and
// decode the file when first run
func LoadConfig(file string) *Config {
	log.Infof("read config from %s", file)
	if _, err := toml.DecodeFile(file, &config); err != nil {
		log.Fatalf("Error: %s", err)
	}
	return config
}

// GetConfig return the config object
func GetConfig() *Config {
	return config
}
