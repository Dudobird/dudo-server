package utils

import (
	"testing"
)

func TestGetFileSizeFromReadable(t *testing.T) {
	testCases := []struct {
		expect uint64
		input  string
	}{
		{
			expect: 1024,
			input:  "1KB",
		},
		{
			expect: 2097152,
			input:  "2.0MB",
		},
		{
			expect: 30 * 1024 * 1024 * 1024,
			input:  "30GB",
		},
		{
			expect: 100 * 1024 * 1024 * 1024 * 1024,
			input:  "100TB",
		},
		{
			expect: 30 * 1024 * 1024 * 1024,
			input:  "30.0GB",
		},
		{
			expect: 0,
			input:  "",
		},
	}

	for _, tc := range testCases {
		Equals(t, int(GetFileSizeFromReadable(tc.input)), int(tc.expect))
	}
}

func TestGetReadableFileSize(t *testing.T) {
	testCases := []struct {
		input  float64
		expect string
	}{
		{
			input:  2000,
			expect: "1.95KB",
		},
		{
			input:  2000000,
			expect: "1.91MB",
		},
		{
			input:  2000000000,
			expect: "1.86GB",
		},
		{
			input:  2000000000000,
			expect: "1.82TB",
		},
		{
			input:  0,
			expect: "0.00B",
		},
		{
			input:  -1,
			expect: "",
		},
	}

	for _, tc := range testCases {
		Equals(t, GetReadableFileSize(tc.input), tc.expect)
	}
}

func TestGenRandomID(t *testing.T) {
	testCases := []struct {
		length       int
		prefix       string
		expectLength int
	}{
		{
			length:       10,
			prefix:       "user",
			expectLength: 15,
		},
		{
			length:       0,
			prefix:       "user",
			expectLength: 4,
		},
		{
			length:       -2,
			prefix:       "user",
			expectLength: 4,
		},
		{
			length:       10,
			prefix:       "",
			expectLength: 10,
		},
	}

	for _, tc := range testCases {
		Equals(t, tc.expectLength, len(GenRandomID(tc.prefix, tc.length)))
	}
}

func TestGetFilePathFolders(t *testing.T) {
	testCases := []struct {
		filePath string
		expect   []string
	}{
		{
			filePath: "/a/b/c/dd/d/1.txt",
			expect:   []string{"a", "b", "c", "dd", "d"},
		},
		{
			filePath: "1.txt",
			expect:   []string{},
		},
		{
			filePath: "",
			expect:   []string{},
		},
		{
			filePath: "/a/b/1.txt",
			expect:   []string{"a", "b"},
		},
		{
			filePath: "1.txt",
			expect:   []string{},
		},
	}

	for _, tc := range testCases {
		Equals(t, tc.expect, GetFilePathFolderList(tc.filePath))
	}
}

func TestGetFileExtention(t *testing.T) {
	testCases := []struct {
		fileName string
		expect   string
	}{
		{
			fileName: "1.doc ",
			expect:   "doc",
		},
		{
			fileName: "abc.doc",
			expect:   "doc",
		},
		{
			fileName: "abc.pdf",
			expect:   "pdf",
		},
		{
			fileName: "abc",
			expect:   "file",
		},
		{
			fileName: "",
			expect:   "",
		},
	}

	for _, tc := range testCases {
		Equals(t, tc.expect, GetFileExtention(tc.fileName))
	}
}
