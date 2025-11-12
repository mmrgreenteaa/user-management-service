package handlers

import (
	"context"
	"errors"
	"log/slog"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/mmrgreenteaa/user-management-service/internal/gen/proto/user_manegement"
	"github.com/mmrgreenteaa/user-management-service/internal/user_management/database/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *UserManagementServer) VerifyCredentials(ctx context.Context, req *pb.UserValidRequest) (*pb.UserValidResponse, error) {

	id, err := s.Db.GetUserId(req.Login, req.Password)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			s.logger.Error("verifyCredentials: get user id invalid login or valid", slog.String("Login", req.Login))
			return nil, status.Error(codes.Unauthenticated, "invalid login or password")
		}
		s.logger.Error("failed get user", slog.String("Error", err.Error()))
		return nil, status.Error(codes.Internal, "")
	}
	s.logger.Info("user id successfully issued")
	return &pb.UserValidResponse{UserId: id}, nil
}

func (s *UserManagementServer) RegistrationUser(ctx context.Context, req *pb.UserRegistrationReq) (*empty.Empty, error) {

	err := s.Db.AddUser(req.Login, req.Password)
	if err != nil {
		s.logger.Error("failed add user in db", slog.String("User login:", req.Login))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return nil, nil
}

func (s *UserManagementServer) EditLogin(ctx context.Context, req *pb.UserLoginEditReq) (*empty.Empty, error) {
	err := s.Db.EditLogin(req.UserId, req.NewLogin)
	if err != nil {
		if errors.Is(err, mongodb.NoRecord) {
			s.logger.Error("edit login no record", slog.String("User id", req.UserId))
			return nil, status.Error(codes.InvalidArgument, "user_id invalid")
		}
		s.logger.Error("falied login edit", slog.String("Error", err.Error()))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return nil, nil
}

func (s *UserManagementServer) DeleteUser(ctx context.Context, req *pb.DeleteReq) (*empty.Empty, error) {

	//Note: добавито чтобы пользователь мог удалить только свои дыннеые
	err := s.Db.DeleteUser(req.UserId)
	if err != nil {
		if errors.Is(err, mongodb.NoRecord) {
			s.logger.Error("faled delete user no record", slog.String("User id", req.UserId))
			return nil, status.Error(codes.InvalidArgument, "user id invalid")
		}

		s.logger.Error("failed delete user", slog.String("Error", err.Error()))
		return nil, status.Error(codes.Internal, err.Error())
	}
	s.logger.Warn("delete user", slog.String("User id", req.UserId))
	return nil, nil
}
