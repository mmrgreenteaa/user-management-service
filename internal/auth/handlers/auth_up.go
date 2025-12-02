// Package handlers implements the gRPC server for the authentication service.
package handlers

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/mmrgreenteaa/user-management-service/internal/auth/datebase/postgresql"
	pb "github.com/mmrgreenteaa/user-management-service/internal/gen/proto/auth"
	"github.com/mmrgreenteaa/user-management-service/internal/gen/proto/user_manegement"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


//  AuthServer represents the gRPC server implementation and holds dependencies like the database.
type AuthServer struct {
	pb.AuthServer
	Db         postgresql.DB
	UserClient user_manegement.UserManagementClient
	logger     *slog.Logger
	cfg        *AuthConfig
}

// AuthConfig holds the configuration settings for the auth service.
type AuthConfig struct {
	Ip         string              `mapstructure:"listen_addr"`
	UserMengIp string              `mapstructure:"user_menedgerIp"`
	DbConfig   postgresql.DbConfig `mapstructure:"posgresql"`
	SecretKey  string
}

// NewAuthServer creates and initializes a new instance of the AuthServer.
func NewAuthServer(db *postgresql.DB, confg *AuthConfig) (*AuthServer, error) {
	conn, err := grpc.NewClient(confg.UserMengIp, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("falied connect to user menedger - %w", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return &AuthServer{
		Db:         *db,
		UserClient: user_manegement.NewUserManagementClient(conn),
		logger:     logger,
		cfg:        confg,
	}, nil
}
