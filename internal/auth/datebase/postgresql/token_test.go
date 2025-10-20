package postgresql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnct(t *testing.T) {
	Сonnect()
}

func TestAddRefresh(t *testing.T) {
	db := Сonnect()

	var tests = []struct {
		name  string
		input struct {
			Token     string
			UserAgent string
			Ip        string
		}
	}{
		{
			name: "Ok test",
			input: struct {
				Token     string
				UserAgent string
				Ip        string
			}{
				"Test",
				"Test",
				"Test",
			},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			err := db.AddRefresh(test.input.Token, test.input.UserAgent, test.input.Ip)
			assert.NoError(t, err)

		})

	}
}
