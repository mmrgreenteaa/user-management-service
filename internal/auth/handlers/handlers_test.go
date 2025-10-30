package handlers

import (
	"context"
	"log"
	"net"
	"testing"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/mmrgreenteaa/user-management-service/internal/auth/datebase/postgresql"
	"github.com/mmrgreenteaa/user-management-service/internal/gen/proto/auth"
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

func setupTestServer() (*grpc.Server, net.Listener) {

	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := &AuthServer{
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

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(CheckJWT)))
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

func TestRefresh(t *testing.T) {

	grpcServer, lis := setupTestServer()
	defer grpcServer.Stop()
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	clinet := genAuth.NewAuthClient(conn)
	require.NoError(t, err)
	tokens, err := CreateRefreshAcsessToken()
	require.NoError(t, err)
	s := &AuthServer{
		Db: *postgresql.Сonnect(),
	}
	s.Db.AddRefresh(tokens.refresh, "Test", "Test")
	md := metadata.Pairs("authorization", tokens.acsess)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	req := genAuth.RefreshRequst{RefreshToken: tokens.refresh}
	res, err := clinet.RefreshToken(ctx, &req)
	assert.NoError(t, err)
	if res.RefreshToken == tokens.refresh {
		assert.Fail(t, "the token has not been updated")
	}
	t.Log(res)
}
func TestLogOut(t *testing.T) {

	grpcServer, lis := setupTestServer()
	defer grpcServer.Stop()
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	clinet := genAuth.NewAuthClient(conn)
	require.NoError(t, err)
	tokens, err := CreateRefreshAcsessToken()
	require.NoError(t, err)
	s := &AuthServer{
		Db: *postgresql.Сonnect(),
	}
	s.Db.AddRefresh(tokens.refresh, "Test", "Test")
	md := metadata.Pairs("authorization", tokens.acsess)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	req := auth.RefreshRequst{RefreshToken: tokens.refresh}
	_, err = clinet.LogOut(ctx, &req)
	assert.NoError(t, err)
}

func TestLogin(t *testing.T) {

	grpcServer, lis := setupTestServer()
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
	log.Println(res)
}
