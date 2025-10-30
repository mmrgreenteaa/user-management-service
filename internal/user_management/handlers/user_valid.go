package handlers

import (
	"context"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/mmrgreenteaa/user-management-service/internal/gen/proto/user_manegement"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *UserManagementServer) VerifyCredentials(ctx context.Context, req *pb.UserValidRequst) (*pb.UserValidResponse, error) {

	id, err := s.DB.GetUserId(req.Login, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &pb.UserValidResponse{UserId: id}, nil
}

func (s *UserManagementServer) RegistrationUser(ctx context.Context, req *pb.UserRegistrationReq) (*empty.Empty, error) {

	err := s.DB.AddUser(req.Login, req.Pass)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}

func (s *UserManagementServer) EditLogin(ctx context.Context, req *pb.UserLoginEditReq) (*empty.Empty, error) {
	err := s.DB.EditLogin(req.OldLogin, req.NewLogin)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return nil, nil
}

func (s *UserManagementServer) DeleteUser(ctx context.Context, req *pb.LoginReq) (*empty.Empty, error) {

	//Note: добавито чтобы пользователь мог удалить только свои дыннеые
	err := s.DB.DeleteUser(req.Login)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return nil, nil
}
