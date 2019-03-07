package storage

import (
	"testing"

	"github.com/Dudobird/dudo-server/utils"
)

func TestGetStorageManager(t *testing.T) {
	storageManager = nil
	utils.Equals(t, (*MinioManager)(nil), GetStorageManager())
}
