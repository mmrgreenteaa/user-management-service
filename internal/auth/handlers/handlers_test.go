package handlers

import (
	"context"
	"log"
	"log/slog"
	"net"
	"os"
	"testing"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/mmrgreenteaa/user-management-service/internal/auth/datebase/postgresql"
	genAuth "github.com/mmrgreenteaa/user-management-service/internal/gen/proto/auth"
	userManegementpb "github.com/mmrgreenteaa/user-management-service/internal/gen/proto/user_manegement"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const userId = "430"

type fakeClinetServ struct {
}

func (f *fakeClinetServ) VerifyCredentials(ctx context.Context, in *userManegementpb.UserValidRequest, opts ...grpc.CallOption) (*userManegementpb.UserValidResponse, error) {
	return &userManegementpb.UserValidResponse{
		UserId: userId,
	}, nil
}
func (f *fakeClinetServ) RegistrationUser(ctx context.Context, in *userManegementpb.UserRegistrationReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, nil
}
func (f *fakeClinetServ) EditLogin(ctx context.Context, in *userManegementpb.UserLoginEditReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, nil
}
func (f *fakeClinetServ) DeleteUser(ctx context.Context, in *userManegementpb.DeleteReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, nil
}

func TestGenerateRefreshToken(t *testing.T) {
	for i := 0; i < 5; i++ {
		ref, err := GenerateRefreshToken()
		assert.NoError(t, err)
		t.Log(ref)
	}

}

type testSetup struct {
	gprcServ *grpc.Server
	lis      net.Listener
	service  *AuthServer
}

func setupTestServer() *testSetup {

	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	fakeClient := &fakeClinetServ{}

	s := &AuthServer{
		Db:         *postgresql.Сonnect(),
		UserClinet: fakeClient,
		logger:     slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	grpcServer := grpc.NewServer()
	genAuth.RegisterAuthServer(grpcServer, s)

	go grpcServer.Serve(lis)
	return &testSetup{
		gprcServ: grpcServer,
		lis:      lis,
		service:  s,
	}
}

func setupTestServerMiddleWare() (*grpc.Server, net.Listener) {

	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	fakeClient := &fakeClinetServ{}
	s := &AuthServer{
		Db:         *postgresql.Сonnect(),
		UserClinet: fakeClient,
		logger:     slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(s.CheckJWT)))
	genAuth.RegisterAuthServer(grpcServer, s)

	go grpcServer.Serve(lis)
	return grpcServer, lis
}

func TestRefresh(t *testing.T) {

	setup := setupTestServer()

	t.Cleanup(func() {
		setup.gprcServ.Stop()
	})

	conn, err := grpc.NewClient(setup.lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() {
		conn.Close()
	})

	clinet := genAuth.NewAuthClient(conn)

	const testLogin = "test-user-for-refresh"
	const testPassword = "password"
	const testUserAgent = "test-agent"
	const testIP = "127.0.0.1"

	loginReq := &genAuth.UserInfoRequest{
		Login:     testLogin,
		Password:  testPassword,
		UserAgent: testUserAgent,
		Ip:        testIP,
	}

	loginRes, err := clinet.Login(context.Background(), loginReq)
	require.NoError(t, err, "Failed to login to get initial tokens")
	require.NotEmpty(t, loginRes.AccessToken)
	require.NotEmpty(t, loginRes.RefreshToken)

	tests := []struct {
		name         string
		accsess      string
		refresh      string
		expectedCode codes.Code
		description  string
	}{
		{
			name:         "Success_ValidToken",
			accsess:      loginRes.AccessToken,
			refresh:      loginRes.RefreshToken,
			expectedCode: codes.OK,
			description:  "A valid refresh token should successfully update tokens",
		},
		{
			name:         "Failure_InvalidToken",
			accsess:      "invalid-access-token",
			refresh:      "invalid-refresh-token",
			expectedCode: codes.Unauthenticated,
			description:  "An invalid refresh token should be rejected",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			ctx := context.WithValue(context.Background(), userIdJwt, userId)
			md := metadata.Pairs("authorization", test.accsess)
			newctx := metadata.NewOutgoingContext(ctx, md)

			req := &genAuth.RefreshRequest{
				RefreshToken: test.refresh,
				UserAgent:    testUserAgent,
				Ip:           testIP,
			}

			resRef, err := clinet.RefreshToken(newctx, req)

			st, ok := status.FromError(err)
			if !ok && err != nil {
				t.Fatalf("Expected a gRPC status error, but got a different error type: %v", err)
			}

			assert.Equal(t, test.expectedCode, st.Code(), test.description)

			if test.expectedCode == codes.OK {
				require.NotNil(t, resRef, "Response should not be nil on success")

				assert.NotEmpty(t, resRef.AccessToken, "New access token should not be empty")
				assert.NotEmpty(t, resRef.RefreshToken, "New refresh token should not be empty")
				assert.NotEqual(t, test.refresh, resRef.RefreshToken, "The refresh token should have been updated")
			}
		})

	}
}

func TestLogOut(t *testing.T) {

	setup := setupTestServer()
	defer setup.gprcServ.Stop()
	conn, err := grpc.NewClient(setup.lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	clinet := genAuth.NewAuthClient(conn)
	require.NoError(t, err)
	tokens, err := CreateRefreshAcsessToken(userId)
	require.NoError(t, err)
	failTokens, err := CreateRefreshAcsessToken("3434")
	require.NoError(t, err)

	tests := []struct {
		name         string
		accsess      string
		refresh      string
		expectedCode codes.Code
		description  string
	}{
		{
			name:         "success_Logout",
			accsess:      tokens.acsess,
			refresh:      tokens.refresh,
			expectedCode: codes.OK,
			description:  "A valid refresh token should successfully update tokens",
		},
		{
			name:         "fail_logout",
			accsess:      failTokens.acsess,
			refresh:      failTokens.refresh,
			expectedCode: codes.Unauthenticated,
			description:  "An invalid refresh token should be rejected",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setup.service.Db.AddRefresh(tokens.refresh, "Test", "Test", "Test")
			md := metadata.Pairs("authorization", test.accsess)
			ctx := metadata.NewOutgoingContext(context.Background(), md)
			req := genAuth.LogoutRequest{
				RefreshToken: test.refresh,
				UserAgent:    "Test",
				Ip:           "Test",
			}
			_, err = clinet.Logout(ctx, &req)
			assert.NoError(t, err)
		})
	}

}

func TestLogin(t *testing.T) {

	setup := setupTestServer()
	defer setup.gprcServ.Stop()
	conn, err := grpc.NewClient(setup.lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	require.NoError(t, err, "")
	defer conn.Close()

	client := genAuth.NewAuthClient(conn)

	req := genAuth.UserInfoRequest{
		Login:     "Test",
		Password:  "Test",
		UserAgent: "Test",
		Ip:        "Test",
	}

	res, err := client.Login(context.Background(), &req)
	assert.NoError(t, err)
	if res.AccessToken == "" {
		assert.Fail(t, "AccessToken the empty")
	}
	reqlog := genAuth.LogoutRequest{
		RefreshToken: res.RefreshToken,
		UserAgent:    "Test",
		Ip:           "Test",
	}
	md := metadata.Pairs("authorization", res.AccessToken)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	_, err = client.Logout(ctx, &reqlog)
	log.Println(res)

}
