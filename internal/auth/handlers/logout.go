package handlers

import (
	"context"

	"github.com/mmrgreenteaa/user-management-service/internal/gen/proto/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthServer) LogOut(ctx context.Context, req *auth.RefreshRequst) (*emptypb.Empty, error) {

	err := s.Db.RefershTokenValid(req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	err = s.Db.DeleteToken(req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return nil, nil
}
