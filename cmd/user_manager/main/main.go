package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/mmrgreenteaa/user-management-service/config"
	usermenpb "github.com/mmrgreenteaa/user-management-service/internal/gen/proto/user_manegement"
	"github.com/mmrgreenteaa/user-management-service/internal/user_management/database/mongodb"
	"github.com/mmrgreenteaa/user-management-service/internal/user_management/handlers"

	"google.golang.org/grpc"
)

func NewGRPCServer() (*grpc.Server, net.Listener) {
	uconf, err := config.GetUsrMnger()
	if err != nil {
		log.Fatal(err)
	}
	
	lis, err := net.Listen("tcp", uconf.Ip)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	db := mongodb.Сonnect(&uconf.DB)
	UserMngServ := handlers.NewUserManagementServer(db)

	grpcServer := grpc.NewServer()
	usermenpb.RegisterUserManagementServer(grpcServer, UserMngServ)
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
