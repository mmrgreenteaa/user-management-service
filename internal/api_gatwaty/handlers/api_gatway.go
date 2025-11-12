package handlers

import (
	"log"
	"log/slog"

	authpb "github.com/mmrgreenteaa/user-management-service/internal/gen/proto/auth"
	usermepb "github.com/mmrgreenteaa/user-management-service/internal/gen/proto/user_manegement"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ApiGatway struct {
	UserMenedger usermepb.UserManagementClient
	AuthServis   authpb.AuthClient
	logger       *slog.Logger
	userid       string
}

func NewApiGatwaty() (*ApiGatway, error) {

	conn, err := grpc.NewClient("127.0.0.1:634840", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("fail to connct ")
	}
	aClient := authpb.NewAuthClient(conn)

	uConn, err := grpc.NewClient("127.0.0.1:63480", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("fail to connct ")
	}
	uClient := usermepb.NewUserManagementClient(uConn)
	return &ApiGatway{
		AuthServis:   aClient,
		UserMenedger: uClient,
	}, nil

}
