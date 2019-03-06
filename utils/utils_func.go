package utils

import "fmt"

// const for storage caculation
const (
	KB = 1024
	MB = KB * 1024
	GB = MB * 1024
	TB = GB * 1024
	PB = TB * 1024
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
