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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	apihandlers "github.com/mmrgreenteaa/user-management-service/internal/api_gatway/handlers"
	usermepb "github.com/mmrgreenteaa/user-management-service/internal/gen/proto/user_manegement"
	"github.com/stretchr/testify/assert"
)

func TestEdit(t *testing.T) {

	gin.SetMode(gin.TestMode)

	ap, err := apihandlers.NewApiGatwaty(&cfg)
	assert.NoError(t, err)

	router := gin.Default()
	router.Use(ap.AuthMiddleware())
	router.POST("/edit_login_user", ap.EditUserLogin)

	type LoginRequest struct {
		//	UserId   string `json:"user_id"`
		NewLogin string `json:"new_login"`
	}

	testSuccess := func(t *testing.T) *usermepb.UserValidResponse {

		req := usermepb.UserValidRequest{
			Login:    loginTst,
			Password: passTst,
		}
		res, err := ap.UserMenedger.VerifyCredentials(context.TODO(), &req)
		assert.NoError(t, err)
		return res

	}(t)

	access := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"user_id": testSuccess.UserId,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	})
	var secretKey = "12333"
	jwtAccess, err := access.SignedString([]byte(secretKey))
	assert.NoError(t, err)

	tests := []struct {
		name     string
		input    *LoginRequest
		statusEx int
	}{
		{

			name: "success",
			input: &LoginRequest{
				//UserId:   testSuccess.UserId,
				NewLogin: "NewtestLogin",
			},
			statusEx: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tt.input)
			assert.NoError(t, err)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/edit_login_user", bytes.NewBuffer(reqBody))
			req.Header.Set("User-Agent", "MyCustomUserAgent/1.0")
			req.RemoteAddr = "192.168.1.100:12345"
			req.Header.Set("Authorization", jwtAccess)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.statusEx, w.Result().StatusCode)
		})

	}

}
func TestDelete(t *testing.T) {

	gin.SetMode(gin.TestMode)

	ap, err := apihandlers.NewApiGatwaty(&cfg)
	assert.NoError(t, err)

	router := gin.Default()
	router.Use(ap.AuthMiddleware())
	router.GET("/delete_user", ap.DeleteAccount)

	testSuccess := func(t *testing.T) *usermepb.UserValidResponse {

		req := usermepb.UserValidRequest{
			Login:    "NewtestLogin",
			Password: passTst,
		}
		res, err := ap.UserMenedger.VerifyCredentials(context.TODO(), &req)
		assert.NoError(t, err)
		return res

	}(t)

	access := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"user_id": testSuccess.UserId,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	})
	var secretKey = "12333"
	jwtAccess, err := access.SignedString([]byte(secretKey))
	assert.NoError(t, err)

	tests := []struct {
		name     string
		statusEx int
	}{
		{
			name:     "success",
			statusEx: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/delete_user", nil)
			req.Header.Set("User-Agent", "MyCustomUserAgent/1.0")
			req.RemoteAddr = "192.168.1.100:12345"
			req.Header.Set("Authorization", jwtAccess)
			router.ServeHTTP(w, req)
			res := w.Result()
			defer res.Body.Close()
			bodyBytes, err := io.ReadAll(res.Body)
			assert.NoError(t, err)
			log.Println(string(bodyBytes))
			assert.Equal(t, tt.statusEx, w.Result().StatusCode)
		})

	}

}
