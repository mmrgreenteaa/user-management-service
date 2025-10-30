package mongodb

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestConnct(t *testing.T) {
	Сonnect()
}

func TestGetUser(t *testing.T) {

	db := Сonnect()

	var tests = []struct {
		name  string
		input struct {
			login string
			pass  string
		}
		err error
	}{
		{
			name: "Ok test",
			input: struct {
				login string
				pass  string
			}{
				"vadim",
				"123",
			},
			err: nil,
		},

		{
			name: "fail test",
			input: struct {
				login string
				pass  string
			}{
				"Test",
				"Test",
			},
			err: mongo.ErrNoDocuments,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			id, err := db.GetUserId(test.input.login, test.input.pass)
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
			log.Println(id)

		})

	}

}

func TestAddUser(t *testing.T) {
	db := Сonnect()
	type testInput struct {
		login string
		pass  string
	}
	tests := []struct {
		name  string
		input testInput
	}{
		{
			name:  "ok",
			input: testInput{login: "vadim", pass: "123"},
		},
		{
			name:  "ok",
			input: testInput{login: "Test", pass: "Test"},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			err := db.AddUser(test.input.login, test.input.pass)
			assert.NoError(t, err)
		})

	}

}

func TestEdit(t *testing.T) {
	db := Сonnect()
	type testInput struct {
		oldLogin string
		newLogin string
	}
	tests := []struct {
		name  string
		input testInput
		err   error
	}{
		{
			name:  "ok",
			input: testInput{oldLogin: "vadim", newLogin: "vadomtest1"},
			err:   nil,
		},
		{
			name:  "old login not found",
			input: testInput{oldLogin: "Testss", newLogin: "vadomtest2"},
			err:   fmt.Errorf("fail: Document not found"),
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			err := db.EditLogin(test.input.oldLogin, test.input.newLogin)
			if test.err != nil {
				if assert.Error(t, err) {
					if err.Error() == test.err.Error() {
						return
					}
				}
			} else {
				assert.NoError(t, err)
			}
		})

	}

}

func TestDeleteUser(t *testing.T) {
	db := Сonnect()

	tests := []struct {
		name  string
		input string
		err   error
	}{
		{
			name:  "ok",
			input: "vadomtest1",
			err:   nil,
		},
		{
			name:  "login not found",
			input: "test",
			err:   fmt.Errorf("fail document not found"),
		},
	}
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			err := db.DeleteUser(test.input)
			if test.err != nil {
				if assert.Error(t, err) {
					if err.Error() == test.err.Error() {
						return
					}
					t.Errorf("expected error %q, got %q", test.err.Error(), err.Error())
				}

			} else {
				assert.NoError(t, err)
			}
		})

	}
}
