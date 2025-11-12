package handlers

import (
	"context"
	"errors"
	"log/slog"

	"github.com/mmrgreenteaa/user-management-service/internal/auth/datebase/postgresql"
	"github.com/mmrgreenteaa/user-management-service/internal/gen/proto/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

// Здесь middleware нужен
func (s *AuthServer) LogOut(ctx context.Context, req *auth.RefreshRequest) (*emptypb.Empty, error) {

	
	userId, ok := ctx.Value(UserIdJwt).(string)
	if !ok {
		s.logger.Error("userid was not found in the context")
		return nil, status.Error(codes.Internal, "failed provided user id")
	}

	idRefresh, err := s.Db.RefershTokenValid(req.RefreshToken, req.UserAgent, req.Ip, userId)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.InvalidArgument, "refresh token not found")
		}
		s.logger.Error("failed refresh valid", "Err:", slog.String("Error", err.Error()))
		return nil, status.Error(codes.Internal, "failed refresh token valid")

	}

	err = s.Db.DeleteToken(idRefresh)
	if err != nil {
		if errors.Is(err, postgresql.NoRecord) {
			s.logger.Warn("falied delete refresh token no record found")
			return nil, status.Error(codes.InvalidArgument, "user id dont't found")
		}
		s.logger.Error("failed delete refresh token", slog.String("Error", err.Error()))
		return nil, status.Error(codes.Internal, "failed delete refresh token")
	}
	return nil, nil
}
