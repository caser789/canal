package mysql

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDefaultError(t *testing.T) {
	tests := []struct {
		name    string
		errCode uint16
		args    []interface{}
		errWant *MySQLError
	}{
		{
			name:    "test has state has format",
			errCode: ER_DUP_KEY,
			args:    []interface{}{"table_a"},
			errWant: &MySQLError{
				Code:    ER_DUP_KEY,
				State:   "23000",
				Message: "Can't write; duplicate key in table 'table_a'",
			},
		},
		{
			name:    "test no state no format",
			errCode: 8888,
			errWant: &MySQLError{
				Code:    8888,
				State:   "HY000",
				Message: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.name)
			errGot := NewDefaultError(tt.errCode, tt.args...)
			assert.Equal(t, tt.errWant, errGot)
		})
	}
}

func TestNewError(t *testing.T) {
	tests := []struct {
		name    string
		errCode uint16
		message string
		errWant *MySQLError
	}{
		{
			name:    "test happy",
			errCode: 8888,
			message: "happy",
			errWant: &MySQLError{
				Code:    8888,
				State:   "HY000",
				Message: "happy",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.name)
			errGot := NewError(tt.errCode, tt.message)
			assert.Equal(t, tt.errWant, errGot)
		})
	}
}
