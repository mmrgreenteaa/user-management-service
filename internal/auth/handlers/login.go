package handlers

import (
	"context"
	"log/slog"

	"github.com/mmrgreenteaa/user-management-service/internal/gen/proto/auth"
	"github.com/mmrgreenteaa/user-management-service/internal/gen/proto/user_manegement"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AuthServer) Login(ctx context.Context, req *auth.UserInfoRequest) (*auth.SignInUserResponse, error) {

	if req.Login == "" || req.Password == "" || req.UserAgent == "" || req.Ip == "" {
		s.logger.Error("empty values requst")
		return nil, status.Error(codes.InvalidArgument, "empty values requst")
	}

	uReq := user_manegement.UserValidRequest{
		Login:    req.Login,
		Password: req.Password,
	}
	res, err := s.UserClinet.VerifyCredentials(ctx, &uReq)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			s.logger.Error("unknow status VerifyCredentials func service user menagement")
			return nil, status.Errorf(codes.Unknown, "user manegement service failed:%v", st.Message())
		}
		s.logger.Error("verifyCredentials fail ", slog.String("Error", err.Error()))
		return nil, status.Errorf(st.Code(), "verify credentials failed: %v", st.Message())
	}

	tokens, err := CreateRefreshAcsessToken(res.UserId)
	if err != nil {
		s.logger.Error("fail create tokens", slog.String("Error", err.Error()))
		return nil, status.Errorf(codes.Internal, "faid create tokens - %v", err.Error())
	}

	err = s.Db.AddRefresh(tokens.refresh, res.UserId, req.UserAgent, req.Ip)
	if err != nil {
		s.logger.Error("fail add refresg in db", slog.String("Error", err.Error()))
		return nil, status.Error(codes.Internal, "failed to store refresh token")
	}

	s.logger.Info("the login was successful")
	return &auth.SignInUserResponse{
		AccessToken:  tokens.acsess,
		RefreshToken: tokens.refresh,
	}, nil

}
