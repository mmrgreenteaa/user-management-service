package main

import (
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-gormigrate/gormigrate/v2"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/mmrgreenteaa/user-management-service/config"
	"github.com/mmrgreenteaa/user-management-service/internal/auth/datebase/postgresql"
	auth "github.com/mmrgreenteaa/user-management-service/internal/auth/handlers"
	migrate "github.com/mmrgreenteaa/user-management-service/internal/auth/migrations"
	genAuth "github.com/mmrgreenteaa/user-management-service/internal/gen/proto/auth"
	"google.golang.org/grpc"
)

var logger = slog.Default()

func NewGRPCServer() (*grpc.Server, net.Listener) {

	authConfg, err := config.GetAuth()
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", authConfg.Ip)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	db := postgresql.Сonnect(&authConfg.DbConfig)

	res := db.Exec("CREATE SCHEMA IF NOT EXISTS tokens_info")
	if res.Error != nil {
		log.Fatal(res.Error)
	}
	m := gormigrate.New(db.DB, gormigrate.DefaultOptions, []*gormigrate.Migration{migrate.M1224CreateRerefreshTokens()})
	err = m.Migrate()
	if err != nil {
		logger.Error("migration error", slog.String("Error", err.Error()))
	}
	as, err := auth.NewAuthServer(db, authConfg)
	if err != nil {
		log.Fatalf("falied new auth server %v", err)
	}
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(as.ParseJWT)))
	genAuth.RegisterAuthServer(grpcServer, as)

	return grpcServer, lis
}

func main() {

	serv, lis := NewGRPCServer()

	go func() {
		log.Printf("gRPC server is starting on %s", lis.Addr())
		if err := serv.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down gRPC server...")
	serv.GracefulStop()
	log.Println("Server gracefully stopped.")

}
