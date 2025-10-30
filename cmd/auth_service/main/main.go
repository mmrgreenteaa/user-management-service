package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/mmrgreenteaa/user-management-service/internal/auth/datebase/postgresql"
	auth "github.com/mmrgreenteaa/user-management-service/internal/auth/handlers"
	genAuth "github.com/mmrgreenteaa/user-management-service/internal/gen/proto/auth"
	"google.golang.org/grpc"
)

func NewGRPCServer() (*grpc.Server, net.Listener) {

	lis, err := net.Listen("tcp", "127.0.0.1:62480")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	db := postgresql.Сonnect()
	as := auth.NewAuthServer(db)
	grpcServer := grpc.NewServer()
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
