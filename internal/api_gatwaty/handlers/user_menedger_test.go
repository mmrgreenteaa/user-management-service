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
			ap.userid = tt.input.UserId
			req.Header.Set("Content-Type", "application/json")
			ctx := context.WithValue(context.Background(), "user_id", tt.input.UserId)
			req.WithContext(ctx)
			router.ServeHTTP(w, req)
			log.Println(w.Header())
			log.Println(w.Body)
			assert.Equal(t, tt.statusEx, w.Code)
		})
	}

}

func TestUserDelete(t *testing.T) {

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
	router.POST("/ping", ap.DeleteAccount)
	type LoginRequest struct {
		Login string `json:"login"`
	}

	tests := []struct {
		name     string
		input    LoginRequest
		statusEx int
	}{
		{

			name: "success",
			input: LoginRequest{
				UserId: "400",
				Login:  "hero2",
			},
			statusEx: http.StatusOK,
		},

		{
			name: "login not corrected",
			input: LoginRequest{
				UserId: "300",
				Login:  "",
			},
			statusEx: http.StatusBadRequest,
		},
	}

	data := LoginRequest{
		Login: "Test",
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/ping", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	log.Println(w.Header())
	log.Println(w.Body)
	assert.Equal(t, 200, w.Code)

}
