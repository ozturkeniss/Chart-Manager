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

func (c *AuthClient) ValidateToken(token string) (interface{}, error) {
	ctx := context.Background()
	req := &pb.ValidateTokenRequest{
		Token: token,
	}

	response, err := c.client.ValidateToken(ctx, req)
	if err != nil {
		return nil, err
	}

	// Convert to map for interface{} compatibility
	result := map[string]interface{}{
		"valid":    response.Valid,
		"user_id":  response.UserId,
		"username": response.Username,
		"role":     response.Role,
	}

	return result, nil
}

func (c *AuthClient) GetUser(userID uint32) (interface{}, error) {
	ctx := context.Background()
	req := &pb.GetUserRequest{
		UserId: userID,
	}

	response, err := c.client.GetUser(ctx, req)
	if err != nil {
		return nil, err
	}

	// Convert to map for interface{} compatibility
	result := map[string]interface{}{
		"success": response.Success,
		"message": response.Message,
	}

	if response.User != nil {
		result["user_id"] = response.User.Id
		result["username"] = response.User.Username
		result["email"] = response.User.Email
		result["role"] = response.User.Role
	}

	return result, nil
}
