package handlers

import (
	"log"
	"log/slog"
	"os"

	"github.com/mmrgreenteaa/user-management-service/internal/auth/datebase/postgresql"
	pb "github.com/mmrgreenteaa/user-management-service/internal/gen/proto/auth"
	"github.com/mmrgreenteaa/user-management-service/internal/gen/proto/user_manegement"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthServer struct {
	pb.AuthServer
	Db         postgresql.DB
	UserClinet user_manegement.UserManagementClient
	logger     *slog.Logger
}

var secretKey = "12333"

func NewAuthServer(db *postgresql.DB) *AuthServer {

	conn, err := grpc.NewClient("127.0.0.1:634840", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("fail to connct ")
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return &AuthServer{
		Db:         *db,
		UserClinet: user_manegement.NewUserManagementClient(conn),
		logger:     logger,
	}
}
