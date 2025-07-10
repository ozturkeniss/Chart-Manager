package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "rancher-manager/api/proto/itemservice"
	"rancher-manager/internal/itemservice/model"
)

type ItemServiceInterface interface {
	GetItem(id string, userID uint32) (*model.Item, error)
	UpdateItem(id string, req *model.UpdateItemRequest, userID uint32) (*model.Item, error)
	DeleteItem(id string, userID uint32) error
}

type ItemGRPCServer struct {
	pb.UnimplementedItemServiceServer
	itemService ItemServiceInterface
}

func NewItemGRPCServer(itemService ItemServiceInterface) *ItemGRPCServer {
	return &ItemGRPCServer{
		itemService: itemService,
	}
}

// authInterceptor extracts user_id from metadata and adds it to context
func authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata not found")
	}

	userIDs := md.Get("user_id")
	if len(userIDs) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "user_id not found in metadata")
	}

	// For simplicity, we'll use the first user_id from metadata
	// In a real implementation, you might want to validate the user_id format
	userIDStr := userIDs[0]

	// Convert string to uint32 (you might want to add proper validation here)
	var userID uint32
	_, err := fmt.Sscanf(userIDStr, "%d", &userID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id format")
	}

	// Add user_id to context
	ctx = context.WithValue(ctx, "user_id", userID)

	return handler(ctx, req)
}

func (s *ItemGRPCServer) GetItem(ctx context.Context, req *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	// Get user from context (set by auth interceptor)
	userID, ok := ctx.Value("user_id").(uint32)
	if !ok {
		return &pb.GetItemResponse{
			Success: false,
			Message: "User not authenticated",
		}, nil
	}

	item, err := s.itemService.GetItem(req.ItemId, userID)
	if err != nil {
		return &pb.GetItemResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.GetItemResponse{
		Success: true,
		Item: &pb.Item{
			Id:          item.ID.Hex(),
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
			Category:    item.Category,
			Stock:       int32(item.Stock),
			CreatedBy:   item.CreatedBy,
			UpdatedBy:   item.UpdatedBy,
			CreatedAt:   item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
		Message: "Item retrieved successfully",
	}, nil
}

func (s *ItemGRPCServer) UpdateStock(ctx context.Context, req *pb.UpdateStockRequest) (*pb.UpdateStockResponse, error) {
	// Get user from context (set by auth interceptor)
	userID, ok := ctx.Value("user_id").(uint32)
	if !ok {
		return &pb.UpdateStockResponse{
			Success: false,
			Message: "User not authenticated",
		}, nil
	}

	updateReq := &model.UpdateItemRequest{
		Stock: int(req.NewStock),
	}

	item, err := s.itemService.UpdateItem(req.ItemId, updateReq, userID)
	if err != nil {
		return &pb.UpdateStockResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.UpdateStockResponse{
		Success: true,
		Item: &pb.Item{
			Id:          item.ID.Hex(),
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
			Category:    item.Category,
			Stock:       int32(item.Stock),
			CreatedBy:   item.CreatedBy,
			UpdatedBy:   item.UpdatedBy,
			CreatedAt:   item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
		Message: "Stock updated successfully",
	}, nil
}

func (s *ItemGRPCServer) DeleteItem(ctx context.Context, req *pb.DeleteItemRequest) (*pb.DeleteItemResponse, error) {
	// Get user from context (set by auth interceptor)
	userID, ok := ctx.Value("user_id").(uint32)
	if !ok {
		return &pb.DeleteItemResponse{
			Success: false,
			Message: "User not authenticated",
		}, nil
	}

	err := s.itemService.DeleteItem(req.ItemId, userID)
	if err != nil {
		return &pb.DeleteItemResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.DeleteItemResponse{
		Success: true,
		Message: "Item deleted successfully",
	}, nil
}

func StartGRPCServer(itemService ItemServiceInterface, port string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor))
	pb.RegisterItemServiceServer(grpcServer, NewItemGRPCServer(itemService))

	fmt.Printf("gRPC Item Server listening on port %s\n", port)
	return grpcServer.Serve(lis)
}
