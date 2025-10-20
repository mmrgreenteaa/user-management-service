package handlers

import (
	"context"
	"log"
	"net"
	"testing"

	authgprc "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/mmrgreenteaa/user-management-service/internal/auth/datebase/postgresql"
	genAuth "github.com/mmrgreenteaa/user-management-service/internal/gen/proto/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

type testServer struct {
	genAuth.UnimplementedTestServiceServer
}

func (s *testServer) EmptyCall(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func TestGenerateRefreshToken(t *testing.T) {
	for i := 0; i < 5; i++ {
		ref, err := GenerateRefreshToken()
		assert.NoError(t, err)
		t.Log(ref)
	}

}

func setupTestServer(t *testing.T) (*grpc.Server, net.Listener) {

	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := &Server{
		Db: *postgresql.Сonnect(),
	}

	grpcServer := grpc.NewServer()
	genAuth.RegisterAuthServer(grpcServer, s)

	go grpcServer.Serve(lis)
	return grpcServer, lis
}

func setupTestServerMiddleWare() (*grpc.Server, net.Listener) {

	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(authgprc.UnaryServerInterceptor(CheckJWT)))
	genAuth.RegisterTestServiceServer(grpcServer, &testServer{})

	go grpcServer.Serve(lis)
	return grpcServer, lis
}

func TestMiddleWare(t *testing.T) {

	grpcServer, lis := setupTestServerMiddleWare()
	defer grpcServer.Stop()
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	require.NoError(t, err, "")
	defer conn.Close()

	clinet := genAuth.NewTestServiceClient(conn)

	tokens, err := CreateRefreshAcsessToken()
	require.NoError(t, err)

	md := metadata.Pairs("authorization", tokens.acsess)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	_, err = clinet.EmptyCall(ctx, &emptypb.Empty{})
	assert.NoError(t, err)

}

func TestLogin(t *testing.T) {

	grpcServer, lis := setupTestServer(t)
	defer grpcServer.Stop()
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	require.NoError(t, err, "")
	defer conn.Close()

	client := genAuth.NewAuthClient(conn)

	req := genAuth.UserInfoRequest{
		Login:    "Test",
		Password: "Test",
	}

	res, err := client.Login(context.Background(), &req)
	assert.NoError(t, err)
	if res.AccessToken == "" {
		assert.Fail(t, "AccessToken the empty")
	}

}
