package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	pb "rancher-manager/api/proto/authservice"
	"rancher-manager/internal/authservice/service"
)

type AuthGRPCServer struct {
	pb.UnimplementedAuthServiceServer
	authService *service.AuthService
}

func NewAuthGRPCServer(authService *service.AuthService) *AuthGRPCServer {
	return &AuthGRPCServer{
		authService: authService,
	}
}

func (s *AuthGRPCServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	claims, err := s.authService.ValidateToken(req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{
			Valid:   false,
			Message: err.Error(),
		}, nil
	}

	return &pb.ValidateTokenResponse{
		Valid:    true,
		UserId:   uint32(claims.UserID),
		Username: claims.Username,
		Role:     claims.Role,
		Message:  "Token is valid",
	}, nil
}

func (s *AuthGRPCServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.authService.GetProfile(uint(req.UserId))
	if err != nil {
		return &pb.GetUserResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.GetUserResponse{
		Success: true,
		User: &pb.User{
			Id:        uint32(user.ID),
			Username:  user.Username,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
			IsActive:  user.IsActive,
		},
		Message: "User found",
	}, nil
}

func StartGRPCServer(authService *service.AuthService, port string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, NewAuthGRPCServer(authService))

	fmt.Printf("gRPC Auth Server listening on port %s\n", port)
	return grpcServer.Serve(lis)
}
