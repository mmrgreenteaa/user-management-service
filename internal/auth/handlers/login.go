package handlers

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mmrgreenteaa/user-management-service/internal/gen/proto/auth"
	"github.com/mmrgreenteaa/user-management-service/internal/gen/proto/user_manegement"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AuthServer) Login(ctx context.Context, req *auth.UserInfoRequest) (*auth.SignInUserResponse, error) {

	if req.Login == "" || req.Password == "" || req.UserAgent == "" || req.Ip == "" {
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
			return nil, status.Errorf(codes.Unknown, "user manegement service failed:%v", st.Message())
		}
		return nil, status.Errorf(st.Code(), "verify credentials failed: %v", st.Message())
	}
	refresh, err := GenerateRefreshToken()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate refresh token")
	}

	err = s.Db.AddRefresh(refresh, req.UserAgent, req.Ip)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to store refresh token")
	}

	access := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"user_id": res.UserId,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	})

	jwtAccess, err := access.SignedString([]byte(secretKey))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed sing access tocken")
	}

	return &auth.SignInUserResponse{
		AccessToken:  jwtAccess,
		RefreshToken: refresh,
	}, nil

}

type tokens struct {
	acsess  string
	refresh string
}
