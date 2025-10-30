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
	"github.com/stretchr/testify/require"
)

func TestLogIn(t *testing.T) {

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
	router.POST("/ping", ap.LogIn)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/ping", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	log.Println(w.Header())
	log.Println(w.Body)
	assert.Equal(t, 200, w.Code)

}

func TestAuthMiddleWare(t *testing.T) {

	ap, err := NewApiGatwaty()
	assert.NoError(t, err)
	router := gin.Default()
	router.Use(ap.AuthMiddleware())
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, nil)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	req.Header.Add("Authorization", "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjE3NDQ0NzR9.xnnXveaPsgFGyW0cV09_n5gfc4aEmBm9XF11z956onMs1yrasqa-Ekdje2kqwLHC7h5xyV8lSrbi97Xocs9ZQA")
	router.ServeHTTP(w, req)
	log.Println(w.Header())

	assert.Equal(t, 200, w.Code)

}

func TestLogOut(t *testing.T) {

	ap, err := NewApiGatwaty()
	require.NoError(t, err)
	router := gin.Default()
	router.GET("/ping", ap.LogOut)
	router.Use(ap.AuthMiddleware())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	req.Header.Add("Authorization", "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjE3NDQ0NzR9.xnnXveaPsgFGyW0cV09_n5gfc4aEmBm9XF11z956onMs1yrasqa-Ekdje2kqwLHC7h5xyV8lSrbi97Xocs9ZQA")

	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "3435621d4246e41f420eaafc38c9ae62e28a9969eb8302464c627074aa366dc0"})
	router.ServeHTTP(w, req)

	log.Println(w.Header())
	log.Println(w.Body)
	assert.Equal(t, 200, w.Code)

}

func TestRefresh(t *testing.T) {

	ap, err := NewApiGatwaty()
	require.NoError(t, err)
	router := gin.Default()
	router.GET("/ping", ap.RefreshToken)
	router.Use(ap.AuthMiddleware())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	req.Header.Add("Authorization", "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjE3NDc5NjV9.DW33V0viO1F15YrKofsC07o_O1RxysT7OrnPRo72gUBKZaVFCURU0ZsutATvlbmQaqDZ8NbZysQDQhbZgH-p_g")

	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "806182157365f22c3a75f941925f174f50ceeced76a3450f3bc3c065ca0043d6"})
	router.ServeHTTP(w, req)

	log.Println(w.Header())
	log.Println(w.Body)
	assert.Equal(t, 200, w.Code)
}

