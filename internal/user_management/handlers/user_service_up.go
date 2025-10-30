package handlers

import (
	pb "github.com/mmrgreenteaa/user-management-service/internal/gen/proto/user_manegement"
	"github.com/mmrgreenteaa/user-management-service/internal/user_management/database/mongodb"
)

type UserManagementServer struct {
	*mongodb.DB
	pb.UserManegementServer
}

func NewUserManagementServer(db *mongodb.DB) *UserManagementServer {

	return &UserManagementServer{
		DB: db,
	}

}
