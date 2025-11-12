package handlers

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"testing"

	pb "github.com/mmrgreenteaa/user-management-service/internal/gen/proto/user_manegement"
	"github.com/mmrgreenteaa/user-management-service/internal/user_management/database/mongodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type testSetup struct {
	lis           net.Listener
	grpcServer    *grpc.Server
	MedgerService *UserManagementServer
}

var id string

const loginT = "test"
const passT = "test"

func setSetupServer() *testSetup {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	s := &UserManagementServer{
		Db:     mongodb.Сonnect(),
		logger: logger,
	}

	grpcServer := grpc.NewServer()

	pb.RegisterUserManagementServer(grpcServer, s)
	go grpcServer.Serve(lis)
	return &testSetup{
		grpcServer:    grpcServer,
		lis:           lis,
		MedgerService: s,
	}
}

func AddTestDate(ts *testSetup) error {

	err := ts.MedgerService.Db.AddUser(loginT, passT)
	if err != nil {
		return err
	}
	id, err = ts.MedgerService.Db.GetUserId(loginT, passT)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTestDate(ts *testSetup) error {
	if id == "" {
		return fmt.Errorf("id ")
	}
	err := ts.MedgerService.Db.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}

func TestGetUserId(t *testing.T) {
	testSetup := setSetupServer()
	defer testSetup.grpcServer.Stop()
	conn, err := grpc.NewClient(testSetup.lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	require.NoError(t, err)
	defer conn.Close()

	clinet := pb.NewUserManagementClient(conn)

	err = AddTestDate(testSetup)
	require.NoError(t, err)
	tests := []struct {
		name  string
		input pb.UserValidRequest
		err   error
	}{
		{
			name: "ok",
			input: pb.UserValidRequest{
				Login:    loginT,
				Password: passT,
			},
			err: nil,
		},
		{
			name: "user not found",
			input: pb.UserValidRequest{
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
	err = DeleteTestDate(testSetup)
	require.NoError(t, err)

}

func TestRegistrationUser(t *testing.T) {
	testSetup := setSetupServer()
	defer testSetup.grpcServer.Stop()
	conn, err := grpc.NewClient(testSetup.lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	require.NoError(t, err)
	defer conn.Close()

	clinet := pb.NewUserManagementClient(conn)

	tests := []struct {
		name  string
		input pb.UserRegistrationReq
		err   error
	}{
		{
			name: "ok",
			input: pb.UserRegistrationReq{
				Login:    loginT,
				Password: passT,
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
	err = DeleteTestDate(testSetup)
	require.NoError(t, err)
}

func TestEditLogin(t *testing.T) {
	testSetup := setSetupServer()
	defer testSetup.grpcServer.Stop()
	conn, err := grpc.NewClient(testSetup.lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	clinet := pb.NewUserManagementClient(conn)

	err = AddTestDate(testSetup)
	require.NoError(t, err)

	tests := []struct {
		name  string
		input pb.UserLoginEditReq
		err   error
	}{
		{
			name: "ok",
			input: pb.UserLoginEditReq{
				UserId:   "30",
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
	err = DeleteTestDate(testSetup)
	require.NoError(t, err)
}

func TestDeleteUser(t *testing.T) {

	testSetup := setSetupServer()
	defer testSetup.grpcServer.Stop()
	conn, err := grpc.NewClient(testSetup.lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	clinet := pb.NewUserManagementClient(conn)

	err = AddTestDate(testSetup)
	require.NoError(t, err)

	tests := []struct {
		name  string
		input pb.DeleteReq
		err   error
	}{
		{
			name: "ok",
			input: pb.DeleteReq{
				UserId: "300",
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
