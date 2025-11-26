package handlers

import (
	"log/slog"
	"os"

	pb "github.com/mmrgreenteaa/user-management-service/internal/gen/proto/user_manegement"
	"github.com/mmrgreenteaa/user-management-service/internal/user_management/database/mongodb"
)

type UserManagementServer struct {
	Db *mongodb.DB
	pb.UserManagementServer
	logger *slog.Logger
}
type UserMengConfig struct {
	Ip   string           `mapstructure:"listen_addr"`
	DB   mongodb.DbConfig `mapstructure:"mongodb"`

}

func NewUserManagementServer(db *mongodb.DB) *UserManagementServer {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	return &UserManagementServer{
		Db:     db,
		logger: logger,
	}

}
