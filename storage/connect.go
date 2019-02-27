package storage

import (
	"github.com/Dudobird/dudo-server/config"
	minio "github.com/minio/minio-go"
	log "github.com/sirupsen/logrus"
)

var handler *minio.Client

// InitConnection create init connection
func InitConnection() (Storage, error) {
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
		log.Errorln(err)
		return nil, err
	}
	log.Infoln("connect object storage success")
	minioManager := NewMinioManager(storageConnect)
	return minioManager, nil
}
