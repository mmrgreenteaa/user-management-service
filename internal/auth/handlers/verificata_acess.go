package handlers

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mmrgreenteaa/user-management-service/internal/gen/proto/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userIdKeyType struct{}

var UserIdJwt userIdKeyType

func (s *AuthServer) VerificatAccess(ctx context.Context, access *auth.AccessRequest) (*emptypb.Empty, error) {
	token, err := jwt.Parse(access.Token, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS512 {
			return nil, fmt.Errorf("wrong method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		s.logger.Error("falied token parse", slog.String("Error:", err.Error()))
		return nil, status.Error(codes.Unauthenticated, "token parse error")
	}

	if !token.Valid {
		s.logger.Error("the token is not valid", slog.String("Access token:", access.Token))
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		s.logger.Error("token claims are of invalid type")
		return nil, status.Error(codes.Internal, "token claims are of invalid type")
	}
	expFloat, ok := claims["exp"].(float64)
	if !ok {
		s.logger.Error("exp claim is missing or not a float value")
		return nil, status.Error(codes.Unauthenticated, "exp claim is missing or not a float value")
	}

	expTime := time.Unix(int64(expFloat), 0)

	if time.Now().After(expTime) {
		s.logger.Error("token isa expired", slog.String("Access token", access.Token))
		return nil, status.Error(codes.Unauthenticated, "token isa expired")
	}
	return nil, nil
}

func (s *AuthServer) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	log.Println("client is calling method:", fullMethodName)
	metods := map[string]bool{
		"/github.com.mmrgreenteaa.user.management.service.interal.auth1.Auth/Login": true,
	}
	if metods[fullMethodName] {
		return ctx, nil
	}
	return s.CheckJWT(ctx)
}

func (s AuthServer) CheckJWT(ctx context.Context) (context.Context, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		s.logger.Error("failed get context:")
		return nil, status.Errorf(codes.Internal, "failed get context:")
	}
	tokens := md.Get("Authorization")

	token, err := jwt.Parse(tokens[0], func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS512 {
			return nil, fmt.Errorf("wrong method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		s.logger.Error("falied token parse", slog.String("Error:", err.Error()))
		return nil, status.Errorf(codes.InvalidArgument, "falied token parse error %v", err)
	}

	if !token.Valid {
		s.logger.Error("the token is not valid", slog.String("Access token:", tokens[0]))
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		s.logger.Error("token claims are of invalid type")
		return nil, status.Error(codes.Internal, "token claims are of invalid type")
	}
	expFloat, ok := claims["exp"].(float64)
	if !ok {
		s.logger.Error("exp claim is missing or not a float value")
		return nil, status.Error(codes.Unauthenticated, "exp claim is missing or not a float value")
	}
	expTime := time.Unix(int64(expFloat), 0)

	requid, ok := claims["user_id"].(string)
	if !ok {
		s.logger.Error("user_id in context not found")
		return nil, status.Error(codes.InvalidArgument, "user_id in context not found")
	}

	if time.Now().After(expTime) {
		s.logger.Error("token isa expired", slog.String("Access token", tokens[0]))
		return nil, status.Error(codes.Unauthenticated, "token isa expired")
	}
	newCtx := context.WithValue(ctx, UserIdJwt, requid)
	return newCtx, nil
}
