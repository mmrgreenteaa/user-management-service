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

	authpb "github.com/mmrgreenteaa/user-management-service/internal/gen/proto/auth"
	mocks "github.com/mmrgreenteaa/user-management-service/internal/mocks"
	"github.com/stretchr/testify/assert"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	loginTest = "testLogin"
	passTest  = "testPass"
	emailTest = "testEmail"
)

func TestLogIn(t *testing.T) {

	authMock := mocks.AuthServiceClientMock{
		LoginFunc: func(ctx context.Context, in *authpb.UserInfoRequest, opts ...grpc.CallOption) (*authpb.SignInUserResponse, error) {

			log.Println("requst varible:", in.Login, in.Password)
			if in.Login != loginTest && in.Password != passTest {
				return nil, status.Errorf(codes.Unauthenticated, "log or pass invalid")
			}
			return &authpb.SignInUserResponse{
				AccessToken:  "acesstoken",
				RefreshToken: "refreshtoke",
			}, nil
		},
	}

	apgm := ApiGatway{
		AuthServis:   &authMock,
		UserMenedger: nil,
		logger:       slog.Default(),
	}

	router := gin.Default()
	router.POST("/ping", apgm.LogIn)

	type LoginRequest struct {
		Login string `json:"login"`
		Pass  string `json:"password"`
	}
	tests := []struct {
		name   string
		input  LoginRequest
		status int
	}{
		{
			name: "success",
			input: LoginRequest{
				Login: loginTest,
				Pass:  passTest,
			},
			status: http.StatusOK,
		},
		{
			name: "the user will not be able to log in",
			input: LoginRequest{
				Login: "blabla",
				Pass:  "blabla",
			},
			status: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			reqBody, err := json.Marshal(tt.input)
			if err != nil {
				panic(err)
			}
			req, err := http.NewRequest("POST", "/ping", bytes.NewBuffer(reqBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)
			log.Println(w.Header())
			log.Println(w.Body)
			assert.Equal(t, tt.status, w.Code)

		})

	}

}

func TestAuthMiddleWare(t *testing.T) {

	authmock := mocks.AuthServiceClientMock{
		VerifyAccessFunc: func(ctx context.Context, in *authpb.AccessRequest, opts ...grpc.CallOption) (*authpb.VerifyAccessResponse, error) {
			return &authpb.VerifyAccessResponse{
				UserId: "400",
				Login:  loginTest,
			}, nil
		},
	}

	ap := ApiGatway{
		AuthServis:   &authmock,
		UserMenedger: nil,
		logger:       slog.Default(),
	}

	type tokensReq struct {
		access  string
		refersh string
	}

	tests := []struct {
		name   string
		input  string
		status int
	}{
		{
			name:   "the data is filled in",
			input:  "123213",
			status: http.StatusOK,
		},
		{
			name:   "the date is not filled in",
			input:  "",
			status: http.StatusUnauthorized,
		},
	}

	router := gin.Default()
	router.Use(ap.AuthMiddleware())
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, nil)
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/ping", nil)
			req.Header.Add("Authorization", tt.input)
			router.ServeHTTP(w, req)
			log.Println(w.Header())

			assert.Equal(t, tt.status, w.Code)

		})
	}

}

func TestLogOut(t *testing.T) {

	authmock := mocks.AuthServiceClientMock{
		LogoutFunc: func(ctx context.Context, in *authpb.LogoutRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
			return &emptypb.Empty{}, nil
		},
	}

	ap := ApiGatway{
		AuthServis:   &authmock,
		UserMenedger: nil,
		logger:       slog.Default(),
	}

	router := gin.Default()
	router.GET("/ping", ap.LogOut)

	tests := []struct {
		name         string
		userId       string
		access       string
		refersh      string
		refCookcheck bool
		status       int
	}{
		{
			name:         "sucess",
			userId:       "400",
			refersh:      "123",
			refCookcheck: true,
			status:       http.StatusOK,
		},
		{
			name:         "incorrect refresh",
			userId:       "123123",
			refersh:      "123",
			refCookcheck: false,
			status:       http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/ping", nil)

			if tt.refCookcheck {
				req.AddCookie(&http.Cookie{Name: "refresh_token", Value: tt.refersh})
			}
			router.ServeHTTP(w, req)

			log.Println(w.Header())
			log.Println(w.Body)
			assert.Equal(t, tt.status, w.Code)
		})
	}

}

func TestRefresh(t *testing.T) {

	authmock := mocks.AuthServiceClientMock{
		RefreshTokenFunc: func(ctx context.Context, in *authpb.RefreshRequest, opts ...grpc.CallOption) (*authpb.SignInUserResponse, error) {
			return &authpb.SignInUserResponse{}, nil
		},
	}
	ap := ApiGatway{
		AuthServis:   &authmock,
		UserMenedger: nil,
		logger:       slog.Default(),
	}
	router := gin.Default()
	router.GET("/ping", ap.RefreshToken)

	tests := []struct {
		name         string
		userId       string
		access       string
		refersh      string
		refCookcheck bool
		status       int
	}{
		{
			name:         "sucess",
			userId:       "400",
			access:       "123",
			refersh:      "123",
			refCookcheck: true,
			status:       http.StatusOK,
		},
		{
			name:         "incorrect refresh and access",
			userId:       "123123",
			access:       "",
			refersh:      "123",
			refCookcheck: false,
			status:       http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/ping", nil)
			w := httptest.NewRecorder()
			req.Header.Add("Authorization", tt.access)
			if tt.refCookcheck {
				req.AddCookie(&http.Cookie{Name: "refresh_token", Value: tt.refersh})
			}
			router.ServeHTTP(w, req)

			log.Println(w.Header())
			log.Println(w.Body)
			assert.Equal(t, tt.status, w.Code)
		})
	}

}
