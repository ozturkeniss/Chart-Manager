package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "rancher-manager/api/proto/authservice"
)

type AuthClient struct {
	client pb.AuthServiceClient
}

func NewAuthClient(authServiceAddr string) (*AuthClient, error) {
	conn, err := grpc.Dial(authServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %v", err)
	}

	client := pb.NewAuthServiceClient(conn)
	return &AuthClient{client: client}, nil
}

func (c *AuthClient) ValidateToken(tokenString string) (*pb.ValidateTokenResponse, error) {
	ctx := context.Background()
	req := &pb.ValidateTokenRequest{
		Token: tokenString,
	}

	return c.client.ValidateToken(ctx, req)
}

func (c *AuthClient) GetUser(userID uint32) (*pb.GetUserResponse, error) {
	ctx := context.Background()
	req := &pb.GetUserRequest{
		UserId: userID,
	}

	return c.client.GetUser(ctx, req)
}
