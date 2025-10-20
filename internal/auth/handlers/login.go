package handlers

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	//authgprc	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth "

	"crypto/rand"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mmrgreenteaa/user-management-service/internal/auth/datebase/postgresql"
	"github.com/mmrgreenteaa/user-management-service/internal/gen/proto/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var secretKey = "12333"

type Server struct {
	auth.AuthServer
	Db postgresql.DB
}

func (s *Server) Login(ctx context.Context, req *auth.UserInfoRequest) (*auth.SignInUserResponse, error) {

	if req.Login == "" || req.Password == "" || req.UserAgent == "" || req.Ip == "" {
		return nil, status.Error(codes.InvalidArgument, "empty values requst")
	}

	//NOTE провекрка на базу в сервесе упрвления пользователями

	refresh, err := GenerateRefreshToken()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to sign access token")
	}

	s.Db.AddRefresh(refresh, req.UserAgent, req.Ip)

	access := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		//NOTE: user_id
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	})

	jwtAccess, err := access.SignedString([]byte(secretKey))
	if err != nil {
		return nil, status.Error(codes.Internal, "gerete acess tocken")
	}
	//NOTE: как лучше всего возврощать как значение или в констекскетьт
	return &auth.SignInUserResponse{
		AccessToken:  jwtAccess,
		RefreshToken: refresh,
	}, nil

}

func CheckJWT(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "fail context:")
	}
	tokens := md.Get("authorization")

	token, err := jwt.Parse(tokens[0], func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS512 {
			return nil, fmt.Errorf("wrong method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "token parse error %v", err)
	}

	if !token.Valid {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "token claims are of invalid type")
	}
	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "exp claim is missing or not a float value")
	}

	expTime := time.Unix(int64(expFloat), 0)

	if time.Now().After(expTime) {
		return nil, status.Error(codes.Unauthenticated, "token has expired")
	}

	//NOTE:создание нового контекста для передачи user id из jwt
	return ctx, nil
}

type tokens struct {
	acsess  string
	refresh string
}

func CreateRefreshAcsessToken() (*tokens, error) {
	refresh, err := GenerateRefreshToken()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to sign access token")
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		//NOTE: user_id
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	})

	jwtAccess, err := access.SignedString([]byte(secretKey))
	if err != nil {
		return nil, status.Error(codes.Internal, "gerete acess tocken")
	}
	return &tokens{jwtAccess, refresh}, nil
}

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
