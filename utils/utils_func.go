package utils

import (
	"crypto/rand"
	"fmt"
	"path/filepath"
	"strings"
)

// const for storage caculation
const (
	KB = 1024
	MB = KB * 1024
	GB = MB * 1024
	TB = GB * 1024
	PB = TB * 1024
)

const (
	idSource      = "0123456789qwertyuioasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
	lenOfIDSource = byte(len(idSource))
)

// GetReadableFileSize get the storage
// size to human readable
func GetReadableFileSize(size float64) string {
	switch {
	case size >= PB:
		return fmt.Sprintf("%.2fPB", size/PB)
	case size >= TB:
		return fmt.Sprintf("%.2fTB", size/TB)
	case size >= GB:
		return fmt.Sprintf("%.2fGB", size/GB)
	case size >= MB:
		return fmt.Sprintf("%.2fMB", size/MB)
	case size >= KB:
		return fmt.Sprintf("%.2fKB", size/KB)
	case size >= 0:
		return fmt.Sprintf("%.2fB", size)
	default:
		return ""
	}
}

// GenRandomID generate a random id with prefix
func GenRandomID(prefix string, length int) string {
	if length <= 0 {
		return prefix
	}
	id := make([]byte, length)
	rand.Read(id)
	for i, b := range id {
		id[i] = idSource[b%lenOfIDSource]
	}
	if prefix == "" {
		return string(id)
	}
	return fmt.Sprintf("%s_%s", prefix, string(id))
}

// GetFileExtention get the extensions from file
// if no extensions just return 'file'
func GetFileExtention(fileName string) string {
	if fileName == "" {
		return ""
	}
	name := strings.TrimSpace(fileName)
	extension := filepath.Ext(name)
	if len(extension) > 0 && extension[0] == '.' {
		return extension[1:]
	}
	return "file"
}
