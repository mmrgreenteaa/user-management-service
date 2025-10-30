package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegistrate(t *testing.T) {
	type LoginRequest struct {
		Login string `json:"login"`
		Pass  string `json:"password"`
	}

	data := LoginRequest{
		Login: "myuser",
		Pass:  "mypass",
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	ap, err := NewApiGatwaty()
	assert.NoError(t, err)
	router := gin.Default()
	router.POST("/ping", ap.RegistrateUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/ping", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	log.Println(w.Header())
	log.Println(w.Body)
	assert.Equal(t, 200, w.Code)

}

func TestLoginEdit(t *testing.T) {
	type LoginRequest struct {
		NewLogin string `json:"old_login"`
		OldLogin string `json:"new_login"`
	}

	data := LoginRequest{
		NewLogin: "myuser",
		OldLogin: "myusertest",
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	ap, err := NewApiGatwaty()
	assert.NoError(t, err)
	router := gin.Default()
	router.POST("/ping", ap.EditUserLogin)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/ping", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	log.Println(w.Header())
	log.Println(w.Body)
	assert.Equal(t, 200, w.Code)

}





func TestUserDelete(t *testing.T) {
	type LoginRequest struct {
		Login string `json:"login"`
	}

	data := LoginRequest{
		Login: "Test",
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	ap, err := NewApiGatwaty()
	assert.NoError(t, err)
	router := gin.Default()
	router.POST("/ping", ap.DeleteAccount)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/ping", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	log.Println(w.Header())
	log.Println(w.Body)
	assert.Equal(t, 200, w.Code)

}
