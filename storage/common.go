package storage

// Storage is a interface for upload and get files
// it could be local storage or object storage like s3
type Storage interface {
	// Upload file to storage
	// filepath is the temp upload file path for upload to storage
	Upload(filePath, fileName, bucket string) (path string, err error)

	// Download file from storage
	// filePath is the temp file path for download from storage
	Download(filePath, fileName, bucket string) error

	// Delete file from storage
	Delete(fileName, bucket string) error

	// CleanBucket remove all files from a bucket
	CleanBucket(bucket string) []error
}
