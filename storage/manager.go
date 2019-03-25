package storage

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/Dudobird/dudo-server/models"
	"github.com/Dudobird/dudo-server/utils"
	minio "github.com/minio/minio-go"
	log "github.com/sirupsen/logrus"
)

// MinioManager manager all process with storage
type MinioManager struct {
	Handler      *minio.Client
	TempFilePath string
	Location     string
}

// NewMinioManager create a new storage manager
func NewMinioManager(handler *minio.Client) *MinioManager {
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

// Upload will upload file
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
	exist, err := m.Handler.BucketExists(bucketName)
	if err != nil {
		return err
	}
	if exist == false {
		return errors.New("bucket not exist")
	}
	return m.Handler.FGetObject(bucketName, fileName, filePath, minio.GetObjectOptions{})
}

// Delete download files from minio
func (m *MinioManager) Delete(fileName string, bucketName string) error {
	exist, err := m.Handler.BucketExists(bucketName)
	if err != nil {
		return err
	}
	if exist == false {
		return errors.New("bucket not exist")
	}
	return m.Handler.RemoveObject(bucketName, fileName)
}

// CleanBucket delete  all files in one bucket from minio
func (m *MinioManager) CleanBucket(bucketName string) []error {
	exist, err := m.Handler.BucketExists(bucketName)
	if err != nil {
		return []error{err}
	}
	if exist == false {
		return []error{errors.New("bucket not exist")}
	}
	objectsForDeleteCh := make(chan string)
	errs := []error{}
	go func() {
		doneCh := make(chan struct{})
		defer close(doneCh)
		isRecursive := true
		objects := m.Handler.ListObjectsV2(bucketName, "", isRecursive, doneCh)

		for object := range objects {
			if object.Err != nil {
				log.Errorf("Delete from bucket %s error : %s", bucketName, object.Err)
				errs = append(errs, object.Err)
				continue
			}
			objectsForDeleteCh <- object.Key
		}
		close(objectsForDeleteCh)
	}()
	errCh := m.Handler.RemoveObjects(bucketName, objectsForDeleteCh)
	for e := range errCh {
		log.Errorf("Delete from bucket %s error : %v", bucketName, e.Err)
		errs = append(errs, e.Err)
	}
	return errs
}

// RemoveBucket  delete a bucket and  all files in one bucket if force = true
// from minio
func (m *MinioManager) RemoveBucket(bucketName string, force bool) error {
	if force == true {
		m.CleanBucket(bucketName)
	}
	return m.Handler.RemoveBucket(bucketName)
}

// DownloadFolder will download all files based on files
func (m *MinioManager) DownloadFolder(tempFolderPath, folderName string, files map[string][]models.StorageFile) (string, []error) {
	errors := []error{}
	folderPath := filepath.Join(tempFolderPath, folderName)
	zipFilePath := folderPath + ".zip"
	filePathList := []string{}
	for path, fileList := range files {
		for _, file := range fileList {
			filePath := filepath.Join(tempFolderPath, path, file.FileName)
			if file.IsDir == true {
				// create folder
				os.MkdirAll(filePath, os.ModePerm)
				continue
			}
			err := m.Download(filePath, file.ID, file.Bucket)
			if err != nil {
				log.Errorf("download file %s error: %s", file.FileName, err)
				errors = append(errors, err)
			}
			filePathList = append(filePathList, filePath)
		}
	}
	if len(errors) > 0 {
		return zipFilePath, errors
	}
	// compress file
	err := utils.ZipFiles(zipFilePath, filePathList, tempFolderPath)
	if err != nil {
		return "", append(errors, err)
	}
	return zipFilePath, nil
}
