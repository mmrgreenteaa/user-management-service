package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type tokens struct {
	acsess  string
	refresh string
}

func CreateRefreshAcsessToken(userid string) (*tokens, error) {
	refresh, err := GenerateRefreshToken()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to sign access token")
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"user_id": userid,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
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
		return "", fmt.Errorf("fail generation refresh token %w", err)
	}
	return hex.EncodeToString(b), nil
}
