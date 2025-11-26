package handlers

import (
	"fmt"

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
}
type ApiGatwayConfig struct {
	Ip         string `mapstructure:"listen_addr"`
	AuthIp     string `mapstructure:"authIp"`
	UserMengIP string `mapstructure:"user_managerIp"`
}

func NewApiGatwaty(confg *ApiGatwayConfig) (*ApiGatway, error) {

	aconn, err := grpc.NewClient(confg.AuthIp, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect auth server: %w", err)
	}
	aClient := authpb.NewAuthClient(aconn)

	uConn, err := grpc.NewClient(confg.UserMengIP, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect user menegment: %w", err)
	}
	uClient := usermepb.NewUserManagementClient(uConn)
	return &ApiGatway{
		AuthServis:   aClient,
		UserMenedger: uClient,
		logger:       slog.Default(),
	}, nil

}
