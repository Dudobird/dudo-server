package utils

import (
	"archive/zip"
	"crypto/rand"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	uuid "github.com/satori/go.uuid"
)

// const for storage caculation
const (
	B  = 1
	KB = B * 1024
	MB = KB * 1024
	GB = MB * 1024
	TB = GB * 1024
	PB = TB * 1024
)

var bytesSize = map[string]uint64{
	"b":  B,
	"kb": KB,
	"mb": MB,
	"gb": GB,
	"tb": TB,
	"pb": PB,
}

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

// GetFileSizeFromReadable get file size from readable string
// for example 1MB, return 1024*1024
func GetFileSizeFromReadable(size string) uint64 {
	size = strings.ToLower(size)
	lastDigit := 0
	for _, r := range size {
		if !(unicode.IsDigit(r) || r == '.') {
			break
		}
		lastDigit++
	}
	num := size[:lastDigit]
	f, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return 0
	}

	extra := strings.ToLower(strings.TrimSpace(size[lastDigit:]))
	if m, ok := bytesSize[extra]; ok {
		f = f * float64(m)
		if f >= math.MaxUint64 {
			return 0
		}
		return uint64(f)
	}

	return 0
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

// GetFilePathFolderList split path into many folders list
// without filename
func GetFilePathFolderList(path string) []string {
	folders := []string{}
	for {
		dir := filepath.Dir(path)
		parent := filepath.Base(dir)

		if parent == string(filepath.Separator) || dir == "." {
			break
		}
		path = dir
		folders = append([]string{parent}, folders...)
	}
	return folders
}

// ZipFiles compresses one or many files into a single zip archive file.
// filename is the output zip file's name.
// files is a list of files to add to the zip.
func ZipFiles(filename string, files []string, topFolderPath string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = addFileToZip(zipWriter, file, topFolderPath); err != nil {
			return err
		}
	}
	return nil
}

func addFileToZip(zipWriter *zip.Writer, filename string, topFolderPath string) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	// add relative path
	header.Name, _ = filepath.Rel(topFolderPath, filename)

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

// ValidateUUID validate a uuid string
// return true when is valid, or false
func ValidateUUID(id string) bool {
	_, err := uuid.FromString(id)
	if err != nil {
		return false
	}
	return true
}
