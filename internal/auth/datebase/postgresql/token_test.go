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

func TestRefershTokenValid(t *testing.T) {
	db := Сonnect()

	type input struct {
		refresh   string
		ip        string
		userAgent string
	}

	var tests = []struct {
		name  string
		input input
		err   error
	}{
		{
			name: "Ok test",
			input: input{
				refresh:   " Test",
				ip:        "Test",
				userAgent: "Test",
			},
			err: nil,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			_, err := db.RefershTokenValid(test.input.refresh, test.input.userAgent, test.input.ip)
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

func TestEditRefresh(t *testing.T) {
	db := Сonnect()

	type input struct {
		oldRef string
		newRef string
	}

	var tests = []struct {
		name  string
		input input
		err   error
	}{
		{
			name: "Ok test",
			input: input{
				oldRef: "Test",
				newRef: "test1",
			},
			err: nil,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			err := db.EditRefreshToken(test.input.oldRef, test.input.newRef)
			if test.err != nil {

				if assert.Error(t, err) {

					assert.ErrorIs(t, err, test.err)
				}
			} else {
				assert.NoError(t, err)
			}
		})

	}

}
