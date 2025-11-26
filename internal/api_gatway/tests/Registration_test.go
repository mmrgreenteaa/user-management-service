package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	apihandlers "github.com/mmrgreenteaa/user-management-service/internal/api_gatway/handlers"
	authpb "github.com/mmrgreenteaa/user-management-service/internal/gen/proto/auth"

	"github.com/stretchr/testify/assert"
)

func TestRegestrateUser(t *testing.T) {

	ap, err := apihandlers.NewApiGatwaty(&cfg)
	assert.NoError(t, err)
	gin.SetMode(gin.DebugMode)

	type ReqReq struct {
		Login string `json:"login"`
		Pass  string `json:"password"`
		Email string `json:"email"`
	}

	tests := []struct {
		name     string
		input    *ReqReq
		exStatus int
	}{
		{
			name: "succsful",
			input: &ReqReq{
				Login: "testLogin",
				Pass:  "testPass",
				Email: "testVadim",
			},
			exStatus: http.StatusOK,
		},
	}
	router := gin.New()
	router.POST("/reg", ap.RegistrateUser)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tt.input)
			assert.NoError(t, err)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/reg", bytes.NewBuffer(reqBody))

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.exStatus, w.Result().StatusCode)
		})

	}

}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.DebugMode)
	ap, err := apihandlers.NewApiGatwaty(&cfg)
	assert.NoError(t, err)

	type Req struct {
		Login string `json:"login"`
		Pass  string `json:"password"`
	}

	tests := []struct {
		name     string
		input    *Req
		exStatus int
	}{
		{
			name: "succsful",
			input: &Req{
				Login: "testLogin",
				Pass:  "testPass",
			},
			exStatus: http.StatusOK,
		},
	}
	router := gin.New()
	router.POST("/login", ap.LogIn)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tt.input)
			assert.NoError(t, err)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
			req.Header.Set("User-Agent", "MyCustomUserAgent/1.0")
			req.RemoteAddr = "192.168.1.100:12345"
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.exStatus, w.Result().StatusCode)
		})

	}

}

func TestRefreshToken(t *testing.T) {
	gin.SetMode(gin.DebugMode)
	ap, err := apihandlers.NewApiGatwaty(&cfg)
	assert.NoError(t, err)

	type LoginReq struct {
		Login string `json:"login"`
		Pass  string `json:"password"`
	}

	tokens := func(t *testing.T) *authpb.SignInUserResponse {

		userReq := authpb.UserInfoRequest{
			Login:     loginTst,
			Password:  passTst,
			Ip:        "192.168.1.100:12345",
			UserAgent: "MyCustomUserAgent/1.0",
		}

		res, err := ap.AuthServis.Login(context.TODO(), &userReq)
		assert.NoError(t, err)
		return res

	}(t)

	tests := []struct {
		name     string
		input    *authpb.SignInUserResponse
		exStatus int
	}{
		{
			name:     "succsful",
			input:    tokens,
			exStatus: http.StatusOK,
		},
	}
	router := gin.New()
	router.GET("/refresh", ap.RefreshToken)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/refresh", nil)
			assert.NoError(t, err)
			req.Header.Set("User-Agent", "MyCustomUserAgent/1.0")
			req.Header.Set("Authorization", tt.input.AccessToken)

			cookie := &http.Cookie{Name: "refresh_token", Value: tt.input.RefreshToken}
			req.AddCookie(cookie)
			req.RemoteAddr = "192.168.1.100:12345"
			router.ServeHTTP(w, req)
			res := w.Result()
			defer res.Body.Close()
			bodyBytes, err := io.ReadAll(res.Body)
			assert.NoError(t, err)

			log.Println(string(bodyBytes))
			assert.Equal(t, tt.exStatus, w.Result().StatusCode)
		})

	}

}

func TestLogout(t *testing.T) {
	gin.SetMode(gin.DebugMode)
	ap, err := apihandlers.NewApiGatwaty(&cfg)
	assert.NoError(t, err)

	type LoginReq struct {
		Login string `json:"login"`
		Pass  string `json:"password"`
	}

	tokens := func(t *testing.T) *authpb.SignInUserResponse {

		userReq := authpb.UserInfoRequest{
			Login:     loginTst,
			Password:  passTst,
			Ip:        "192.168.1.100:12345",
			UserAgent: "MyCustomUserAgent/1.0",
		}

		res, err := ap.AuthServis.Login(context.TODO(), &userReq)
		assert.NoError(t, err)
		return res

	}(t)

	tests := []struct {
		name     string
		input    *authpb.SignInUserResponse
		exStatus int
	}{
		{
			name:     "succsful",
			input:    tokens,
			exStatus: http.StatusOK,
		},
	}
	router := gin.New()
	router.GET("/logout", ap.LogOut)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/logout", nil)
			assert.NoError(t, err)
			req.Header.Set("User-Agent", "MyCustomUserAgent/1.0")
			req.Header.Set("Authorization", tt.input.AccessToken)

			cookie := &http.Cookie{Name: "refresh_token", Value: tt.input.RefreshToken}
			req.AddCookie(cookie)
			req.RemoteAddr = "192.168.1.100:12345"
			router.ServeHTTP(w, req)
			res := w.Result()
			defer res.Body.Close()
			bodyBytes, err := io.ReadAll(res.Body)
			assert.NoError(t, err)

			log.Println(string(bodyBytes))
			assert.Equal(t, tt.exStatus, w.Result().StatusCode)
		})

	}

}
