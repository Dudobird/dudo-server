package models

import (
	"testing"

	"github.com/Dudobird/dudo-server/utils"
)

func TestValidate(t *testing.T) {
	type testsuit struct {
		expect bool
		input  *Account
	}
	tests := []testsuit{
		testsuit{
			input:  &Account{},
			expect: false,
		},
		testsuit{
			input:  &Account{Email: "not valid email", Password: "exist"},
			expect: false,
		},
		testsuit{
			input:  &Account{Email: "not valid email", Password: "exist"},
			expect: false,
		},
		testsuit{
			input:  &Account{Email: "", Password: "exist"},
			expect: false,
		},
		testsuit{
			input:  &Account{Email: "test@example.com", Password: "exist"},
			expect: true,
		},
	}
	for _, test := range tests {
		status, _ := test.input.Validate()
		utils.Equals(t, test.expect, status)
	}

}
