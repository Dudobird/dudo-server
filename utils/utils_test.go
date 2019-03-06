package utils

import (
	"testing"
)

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
