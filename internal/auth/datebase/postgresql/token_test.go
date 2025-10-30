package postgresql

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
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

func TestRefershTokenValid(t *testing.T) {
	db := Сonnect()
	var tests = []struct {
		name  string
		input string
		err   error
	}{
		{
			name:  "Ok test",
			input: "Test",
			err:   nil,
		},
		{
			name:  "fail test",
			input: "Tests",
			err:   gorm.ErrRecordNotFound,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			err := db.RefershTokenValid(test.input)
			if test.err != nil {
				// Ожидаем ошибку
				if assert.Error(t, err) {
					// Проверяем, что полученная ошибка соответствует ожидаемой
					assert.ErrorIs(t, err, test.err)
				}
			} else {
				// Ошибок не ожидаем
				assert.NoError(t, err)
			}
		})

	}
}
