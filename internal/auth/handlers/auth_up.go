package handlers

import (
	"log"

	"github.com/mmrgreenteaa/user-management-service/internal/auth/datebase/postgresql"
	"github.com/mmrgreenteaa/user-management-service/internal/gen/proto/auth"
	"github.com/mmrgreenteaa/user-management-service/internal/gen/proto/user_manegement"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthServer struct {
	auth.AuthServer
	Db         postgresql.DB
	UserClinet user_manegement.UserManagementClient
}

var secretKey = "12333"

func NewAuthServer(db *postgresql.DB) *AuthServer {

	conn, err := grpc.NewClient("127.0.0.1:634840", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("fail to connct ")
	}
	return &AuthServer{
		Db:         *db,
		UserClinet: user_manegement.NewUserManagementClient(conn),
	}
}
