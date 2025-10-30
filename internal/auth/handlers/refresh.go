package handlers

import (
	"context"

	"github.com/mmrgreenteaa/user-management-service/internal/gen/proto/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AuthServer) RefreshToken(ctx context.Context, refresh *auth.RefreshRequst) (*auth.SignInUserResponse, error) {

	err := s.Db.RefershTokenValid(refresh.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	tokens, err := CreateRefreshAcsessToken()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	s.Db.AddRefresh(tokens.refresh, "", "")

	return &auth.SignInUserResponse{
		AccessToken:  tokens.acsess,
		RefreshToken: tokens.refresh,
	}, nil
}
