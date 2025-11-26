package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mmrgreenteaa/user-management-service/internal/gen/proto/user_manegement"
	mocks "github.com/mmrgreenteaa/user-management-service/internal/mocks"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestRegistrate(t *testing.T) {

	umm := mocks.UserManagementClientMock{
		RegistrationUserFunc: func(ctx context.Context, in *user_manegement.UserRegistrationReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
			return nil, nil
		},
	}

	ap := ApiGatway{
		UserMenedger: &umm,
		AuthServis:   nil,
		logger:       slog.Default(),
	}

	type LoginRequest struct {
		Login string `json:"login"`
		Pass  string `json:"password"`
		Email string `json:"email"`
	}

	router := gin.Default()
	router.POST("/ping", ap.RegistrateUser)

	tests := []struct {
		name     string
		input    LoginRequest
		statusEx int
	}{
		{

			name: "success",
			input: LoginRequest{
				Login: loginTest,
				Pass:  passTest,
				Email: emailTest,
			},
			statusEx: http.StatusOK,
		},

		{
			name: "login not corrected",
			input: LoginRequest{
				Login: "",
				Pass:  passTest,
				Email: emailTest,
			},
			statusEx: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tt.input)
			assert.NoError(t, err)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/ping", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)
			log.Println(w.Header())
			log.Println(w.Body)
			assert.Equal(t, tt.statusEx, w.Code)
		})
	}
}

func TestLoginEdit(t *testing.T) {

	gin.SetMode(gin.TestMode)
	umm := mocks.UserManagementClientMock{
		EditLoginFunc: func(ctx context.Context, in *user_manegement.UserLoginEditReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
			return nil, nil
		},
	}

	ap := ApiGatway{
		UserMenedger: &umm,
		AuthServis:   nil,
		logger:       slog.Default(),
	}

	router := gin.Default()
	router.POST("/ping", ap.EditUserLogin)

	type LoginRequest struct {
		UserId   string `json:"user_id"`
		NewLogin string `json:"new_login"`
	}

	tests := []struct {
		name     string
		input    LoginRequest
		statusEx int
	}{
		{

			name: "success",
			input: LoginRequest{
				UserId:   "400",
				NewLogin: "hero2",
			},
			statusEx: http.StatusOK,
		},

		{
			name: "login not corrected",
			input: LoginRequest{
				UserId:   "300",
				NewLogin: "",
			},
			statusEx: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tt.input)
			assert.NoError(t, err)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/ping", bytes.NewBuffer(reqBody))
			ctx := gin.CreateTestContextOnly(w, router)

			ctx.Request = req
			ctx.Set("user_id", "400")
			req.Header.Set("Content-Type", "application/json")
			log.Println(ctx)
			ap.EditUserLogin(ctx)
			log.Println(w.Header())
			log.Println(w.Body)
			assert.Equal(t, tt.statusEx, w.Code)
		})
	}

}

func TestUserDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	umm := mocks.UserManagementClientMock{
		DeleteUserFunc: func(ctx context.Context, in *user_manegement.DeleteReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
			return nil, nil
		},
	}
	ap := ApiGatway{
		UserMenedger: &umm,
		AuthServis:   nil,
		logger:       slog.Default(),
	}

	router := gin.Default()
	router.GET("/ping", ap.DeleteAccount)
	type LoginRequest struct {
		Login string `json:"login"`
	}

	tests := []struct {
		name     string
		userId   string
		statusEx int
	}{
		{

			name:     "success",
			userId:   "400",
			statusEx: http.StatusOK,
		},

		{
			name:     "user_id not corrected",
			userId:   "",
			statusEx: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
			ctx := gin.CreateTestContextOnly(w, router)

			ctx.Request = req
			ctx.Set("user_id", tt.userId)
			req.Header.Set("Content-Type", "application/json")
			ap.DeleteAccount(ctx)
			log.Println(w.Header())
			log.Println(w.Body)
			assert.Equal(t, tt.statusEx, w.Code)
		})
	}

}
