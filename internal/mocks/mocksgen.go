package mocks

import (
	authpb "github.com/mmrgreenteaa/user-management-service/internal/gen/proto/auth"
	usermepb "github.com/mmrgreenteaa/user-management-service/internal/gen/proto/user_manegement"
)

type AuthServiceClient interface {
	authpb.AuthClient
}

type UserManagementClient interface {
	usermepb.UserManagementClient
}
