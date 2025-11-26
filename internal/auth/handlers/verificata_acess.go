package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mmrgreenteaa/user-management-service/internal/gen/proto/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"log/slog"
	"strings"
)

type userIdKeyType struct{}

var UserIdJwt userIdKeyType

func (s *AuthServer) VerifyAccess(ctx context.Context, access *auth.AccessRequest) (*auth.VerifyAccessResponse, error) {
	s.logger.Info("requst access token", slog.String("Access", access.Token))
	token, err := jwt.Parse(access.Token, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS512 {
			return nil, fmt.Errorf("wrong method")
		}
		return []byte(s.cfg.SecretKey), nil
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

	uid, ok := claims["user_id"].(string)
	if !ok {
		s.logger.Error("exp claim is missing or not a float value")
		return nil, status.Error(codes.Unauthenticated, "exp claim is missing or not a float value")
	}

	return &auth.VerifyAccessResponse{
		UserId: uid,
	}, nil
}

func (s *AuthServer) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	log.Println("client is calling method:", fullMethodName)
	metods := map[string]bool{
		"/github.com.mmrgreenteaa.user.management.service.interal.auth1.Auth/Login":        true,
		"/github.com.mmrgreenteaa.user.management.service.interal.auth1.Auth/LogOut":       true,
		"/github.com.mmrgreenteaa.user.management.service.interal.auth1.Auth/VerifyAccess": true,
	}

	if metods[fullMethodName] {
		return ctx, nil
	}
	return s.ParseJWT(ctx)
}

func (s AuthServer) ParseJWT(ctx context.Context) (context.Context, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		s.logger.Error("failed get context:")
		return nil, status.Errorf(codes.Internal, "failed get context:")
	}
	tokens := md.Get("Authorization")
	tokenStr := strings.TrimPrefix(tokens[0], "Bearer ")
	log.Println(tokenStr)

	parser := jwt.NewParser(jwt.WithoutClaimsValidation())

	token, err := parser.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS512 {
			return nil, fmt.Errorf("wrong method")
		}
		return []byte(s.cfg.SecretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {

		} else {
			s.logger.Error("falied token parse", slog.String("Error:", err.Error()))
			return nil, status.Errorf(codes.InvalidArgument, "falied token parse error %v", err)
		}
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		s.logger.Error("token claims are of invalid type")
		return nil, status.Error(codes.Internal, "token claims are of invalid type")
	}

	requid, ok := claims["user_id"].(string)
	if !ok {
		s.logger.Error("user_id in context not found")
		return nil, status.Error(codes.InvalidArgument, "user_id in context not found")
	}

	newCtx := context.WithValue(ctx, "user_id", requid)
	//	newCtx := metadata.AppendToOutgoingContext(ctx, "user_id", requid)
	return newCtx, nil
}
