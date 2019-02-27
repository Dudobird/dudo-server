package storage

import (
	minio "github.com/minio/minio-go"
)

// MinioManager manager all process with storage
type MinioManager struct {
	Handler      *minio.Client
	TempFilePath string
	Location     string
}

// NewMinioManager create a new storage manager
func NewMinioManager(handler *minio.Client) Storage {
	if handler == nil {
		return nil
	}
	// Start a goroutine to check and reconnect

	// Return manager object
	return &MinioManager{
		Handler:      handler,
		TempFilePath: "./miniotemp",
		Location:     "dudo",
	}
}

func (m *MinioManager) checkOrCreateBucket(bucket string) error {
	exist, err := m.Handler.BucketExists(bucket)
	if err != nil {
		return err
	}
	if err == nil && !exist {
		err = m.Handler.MakeBucket(bucket, m.Location)
		if err != nil {
			return err
		}
	}
	return nil
}

// Upload will upload file and return a uuid and path
func (m *MinioManager) Upload(filePath string, fileName string, bucketName string) (path string, err error) {
	err = m.checkOrCreateBucket(bucketName)
	if err != nil {
		return
	}
	_, err = m.Handler.FPutObject(bucketName, fileName, filePath, minio.PutObjectOptions{})
	if err != nil {
		return
	}
	return
}

// Download download files from minio
func (m *MinioManager) Download(filePath string, fileName string, bucketName string) error {
	return m.Handler.FGetObject(bucketName, fileName, filePath, minio.GetObjectOptions{})
}

// Delete download files from minio
func (m *MinioManager) Delete(fileName string, bucketName string) error {
	return m.Handler.RemoveObject(bucketName, fileName)
}
