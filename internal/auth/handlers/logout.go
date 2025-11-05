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

func (s *AuthServer) LogOut(ctx context.Context, req *auth.RefreshRequest) (*emptypb.Empty, error) {

	userId, ok := ctx.Value(userIdJwt).(string)
	if !ok {
		s.logger.Error("userid was not found in the context")
		return nil, status.Error(codes.Internal, "failed provided user id")
	}

	id, err := s.Db.RefershTokenValid(req.RefreshToken, req.UserAgent, req.Ip, userId)
	if err != nil {
		if err != nil {
			if errors.Is(err, postgresql.ErrUserAgentMismatch) {
				s.logger.Warn("user agent does not match")
				err = s.Db.DeleteToken(id)
				if err != nil {
					s.logger.Error("failed delete refresh token", slog.String("Error", err.Error()))
					return nil, status.Error(codes.Internal, "failed delete refresh token")
				}
				s.logger.Info("refresh token was delete", "Record id", id)
				return nil, status.Error(codes.Unauthenticated, "fail user agent invalid")
			}
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, status.Error(codes.Unauthenticated, "refresh token does not exist")
			}
			s.logger.Error("failed refresh valid", "Err:", slog.String("Error", err.Error()))
			return nil, status.Error(codes.Internal, "failed refresh token valid")
		}
	}

	err = s.Db.DeleteToken(userId)
	if err != nil {
		s.logger.Error("failed delete refresh token", slog.String("Error", err.Error()))
		return nil, status.Error(codes.Internal, "failed delete refresh token")
	}
	return nil, nil
}
