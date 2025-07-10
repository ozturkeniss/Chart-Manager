package grpc

import (
	"context"
	"fmt"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "rancher-manager/api/proto/itemservice"
)

type ItemClient struct {
	client pb.ItemServiceClient
}

func NewItemClient(itemServiceAddr string) (*ItemClient, error) {
	conn, err := grpc.Dial(itemServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to item service: %v", err)
	}

	client := pb.NewItemServiceClient(conn)
	return &ItemClient{client: client}, nil
}

func (c *ItemClient) GetItem(itemID string, userID uint32) (*pb.GetItemResponse, error) {
	ctx := context.Background()
	md := metadata.New(map[string]string{
		"user_id": strconv.FormatUint(uint64(userID), 10),
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	req := &pb.GetItemRequest{
		ItemId: itemID,
	}

	return c.client.GetItem(ctx, req)
}

func (c *ItemClient) UpdateStock(itemID string, newStock int32, userID uint32) (*pb.UpdateStockResponse, error) {
	ctx := context.Background()
	md := metadata.New(map[string]string{
		"user_id": strconv.FormatUint(uint64(userID), 10),
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	req := &pb.UpdateStockRequest{
		ItemId:   itemID,
		NewStock: newStock,
	}

	return c.client.UpdateStock(ctx, req)
}

func (c *ItemClient) DeleteItem(itemID string, userID uint32) (*pb.DeleteItemResponse, error) {
	ctx := context.Background()
	md := metadata.New(map[string]string{
		"user_id": strconv.FormatUint(uint64(userID), 10),
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	req := &pb.DeleteItemRequest{
		ItemId: itemID,
	}

	return c.client.DeleteItem(ctx, req)
}
