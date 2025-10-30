package handlers

import (
	"context"
	"log"
	"net"
	"testing"

	pb "github.com/mmrgreenteaa/user-management-service/internal/gen/proto/user_manegement"
	"github.com/mmrgreenteaa/user-management-service/internal/user_management/database/mongodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func setSetupServer() (*grpc.Server, net.Listener) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := &UserManagementServer{
		DB: mongodb.Сonnect(),
	}

	grpcServer := grpc.NewServer()

	pb.RegisterUserManegementServer(grpcServer, s)
	go grpcServer.Serve(lis)
	return grpcServer, lis
}

func TestGetUserId(t *testing.T) {
	grpcServer, lis := setSetupServer()
	defer grpcServer.Stop()
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	require.NoError(t, err)
	defer conn.Close()

	clinet := pb.NewUserManegementClient(conn)

	tests := []struct {
		name  string
		input pb.UserValidRequst
		err   error
	}{
		{
			name: "ok",
			input: pb.UserValidRequst{
				Login:    "vadim",
				Password: "123",
			},
			err: nil,
		},
		{
			name: "user not found",
			input: pb.UserValidRequst{
				Login:    "Testd",
				Password: "Test",
			},
			err: mongo.ErrNoDocuments,
		},
	}
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			res, err := clinet.VerifyCredentials(context.Background(), &test.input)
			if test.err != nil {
				assert.Error(t, err)
				t.Log(err)

			} else {
				assert.NoError(t, err)
			}
			t.Logf("user id = %+v", res)
		})

	}

}

func TestRegistrationUser(t *testing.T) {
	grpcServer, lis := setSetupServer()
	defer grpcServer.Stop()
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	require.NoError(t, err)
	defer conn.Close()

	clinet := pb.NewUserManegementClient(conn)

	tests := []struct {
		name  string
		input pb.UserRegistrationReq
		err   error
	}{
		{
			name: "ok",
			input: pb.UserRegistrationReq{
				Login: "vadim2",
				Pass:  "vadim2",
			},
			err: nil,
		},
	}
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			_, err := clinet.RegistrationUser(context.Background(), &test.input)
			if test.err != nil {
				assert.Error(t, err)
				t.Log(err)

			} else {
				assert.NoError(t, err)
			}

		})

	}

}

func TestEditLogin(t *testing.T) {
	grpcServer, lis := setSetupServer()
	defer grpcServer.Stop()
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	clinet := pb.NewUserManegementClient(conn)

	tests := []struct {
		name  string
		input pb.UserLoginEditReq
		err   error
	}{
		{
			name: "ok",
			input: pb.UserLoginEditReq{
				OldLogin: "vadim2",
				NewLogin: "vadim3",
			},
			err: nil,
		},
	}
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			_, err := clinet.EditLogin(context.Background(), &test.input)
			if test.err != nil {
				assert.Error(t, err)
				t.Log(err)

			} else {
				assert.NoError(t, err)
			}

		})

	}

}

func TestDeleteUser(t *testing.T) {

	grpcServer, lis := setSetupServer()
	defer grpcServer.Stop()
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	clinet := pb.NewUserManegementClient(conn)

	tests := []struct {
		name  string
		input pb.LoginReq
		err   error
	}{
		{
			name: "ok",
			input: pb.LoginReq{
				Login: "vadim",
			},
			err: nil,
		},
	}
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			_, err := clinet.DeleteUser(context.Background(), &test.input)
			if test.err != nil {
				assert.Error(t, err)
				t.Log(err)

			} else {
				assert.NoError(t, err)
			}

		})

	}

}
