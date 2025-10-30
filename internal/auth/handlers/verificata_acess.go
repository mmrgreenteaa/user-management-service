package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/mmrgreenteaa/user-management-service/internal/gen/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthServer) VerificatAccess(ctx context.Context, access *auth.AccessRequest) (*emptypb.Empty, error) {
	token, err := jwt.Parse(access.Token, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS512 {
			return nil, fmt.Errorf("wrong method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "token parse error %v", err)
	}

	if !token.Valid {
		return nil, status.Error(codes.InvalidArgument, "invalid token")
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
		return nil, status.Error(codes.Unauthenticated, "token isa expired")
	}
	return nil, nil
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

func SelectiveAuthInterceptor(methods map[string]bool) grpc.UnaryServerInterceptor {
	baseInterceptor := grpc_auth.UnaryServerInterceptor(CheckJWT)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		if methods[info.FullMethod] {
			return baseInterceptor(ctx, req, info, handler)
		}
		// Иначе просто передаём дальше без авторизации
		return handler(ctx, req)
	}
}
