package storage

import (
	"github.com/Dudobird/dudo-server/config"
	minio "github.com/minio/minio-go"
	log "github.com/sirupsen/logrus"
)

var storageManager *MinioManager

// InitStorageManager create storage manager
func InitStorageManager() Storage {
	var err error
	c := config.GetConfig()
	server := c.Storage.Server
	port := c.Storage.Port
	accessKey := c.Storage.AccessKey
	secretKey := c.Storage.SecretKey
	useSSL := false
	log.Infof("try to connect object storage : %s:%s", server, port)
	storageConnect, err := minio.New(server+":"+port, accessKey, secretKey, useSSL)
	if err != nil {
		log.Panicf("connect storage fail: %s", err)
	}
	log.Infoln("connect object storage success")
	storageManager = NewMinioManager(storageConnect)
	return storageManager
}

// GetStorageManager get storage manager object
func GetStorageManager() *MinioManager {
	return storageManager
}
