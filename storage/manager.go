package storage

import minio "github.com/minio/minio-go"

// Manager manager all process with storage
type Manager struct {
	Handler *minio.Client
}

// NewManager create a new storage manager
func NewManager() *Manager {
	handler := GetStorageHandler()
	if handler == nil {
		return nil
	}
	// Start a goroutine to check and reconnect

	// Return manager object
	return &Manager{
		Handler: handler,
	}
}
