package utils

import "testing"

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
