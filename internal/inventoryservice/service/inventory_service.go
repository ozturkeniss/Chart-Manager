package service

import (
	"errors"
	"fmt"

	"rancher-manager/internal/inventoryservice/grpc"
	"rancher-manager/internal/inventoryservice/model"
	"rancher-manager/internal/inventoryservice/repository"
	"rancher-manager/kafka"
)

type InventoryService struct {
	inventoryRepo *repository.InventoryRepository
	authClient    *grpc.AuthClient
	itemClient    *grpc.ItemClient
	publisher     *kafka.Publisher
}

func NewInventoryService(
	inventoryRepo *repository.InventoryRepository,
	authClient *grpc.AuthClient,
	itemClient *grpc.ItemClient,
	publisher *kafka.Publisher,
) *InventoryService {
	return &InventoryService{
		inventoryRepo: inventoryRepo,
		authClient:    authClient,
		itemClient:    itemClient,
		publisher:     publisher,
	}
}

func (s *InventoryService) UpdateStock(itemID string, newStock int, userID uint32) (*model.Inventory, error) {
	// Validate user exists via gRPC
	_, err := s.authClient.GetUser(userID)
	if err != nil {
		return nil, errors.New("unauthorized: invalid user")
	}

	// Get item from ItemService via gRPC
	itemResponse, err := s.itemClient.GetItem(itemID, userID)
	if err != nil {
		return nil, errors.New("item not found")
	}

	if !itemResponse.Success {
		return nil, errors.New("item not found")
	}

	// Check if inventory record exists
	exists, err := s.inventoryRepo.ExistsByItemID(itemID)
	if err != nil {
		return nil, err
	}

	var inventory *model.Inventory
	if !exists {
		// Create new inventory record
		inventory = &model.Inventory{
			ItemID:    itemID,
			Stock:     newStock,
			MinStock:  0,
			MaxStock:  1000,
			UpdatedBy: userID,
		}
		err = s.inventoryRepo.Create(inventory)
	} else {
		// Update existing inventory record
		inventory, err = s.inventoryRepo.GetByItemID(itemID)
		if err != nil {
			return nil, err
		}
		inventory.Stock = newStock
		inventory.UpdatedBy = userID
		err = s.inventoryRepo.Update(inventory)
	}

	if err != nil {
		return nil, err
	}

	// Update stock in ItemService via gRPC
	_, err = s.itemClient.UpdateStock(itemID, int32(newStock), userID)
	if err != nil {
		return nil, fmt.Errorf("failed to update stock in item service: %v", err)
	}

	// Publish Kafka event
	if s.publisher != nil {
		event := &kafka.StockUpdateEvent{
			ItemID:   itemID,
			NewStock: newStock,
			UserID:   userID,
		}
		if err := s.publisher.PublishStockUpdate(event); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Failed to publish stock update event: %v\n", err)
		}
	}

	return inventory, nil
}

func (s *InventoryService) DeleteItem(itemID string, userID uint32) error {
	// Validate user exists via gRPC
	_, err := s.authClient.GetUser(userID)
	if err != nil {
		return errors.New("unauthorized: invalid user")
	}

	// Check if inventory record exists
	exists, err := s.inventoryRepo.ExistsByItemID(itemID)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("inventory record not found")
	}

	// Delete from ItemService via gRPC
	_, err = s.itemClient.DeleteItem(itemID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete item from item service: %v", err)
	}

	// Delete inventory record
	err = s.inventoryRepo.Delete(itemID)
	if err != nil {
		return err
	}

	// Publish Kafka event
	if s.publisher != nil {
		event := &kafka.ItemDeleteEvent{
			ItemID: itemID,
			UserID: userID,
		}
		if err := s.publisher.PublishItemDelete(event); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Failed to publish item delete event: %v\n", err)
		}
	}

	return nil
}

func (s *InventoryService) GetStock(itemID string, userID uint32) (*model.Inventory, error) {
	// Validate user exists via gRPC
	_, err := s.authClient.GetUser(userID)
	if err != nil {
		return nil, errors.New("unauthorized: invalid user")
	}

	inventory, err := s.inventoryRepo.GetByItemID(itemID)
	if err != nil {
		return nil, errors.New("inventory record not found")
	}

	return inventory, nil
}

func (s *InventoryService) GetAllItems(userID uint32) ([]*model.Inventory, error) {
	// Validate user exists via gRPC
	_, err := s.authClient.GetUser(userID)
	if err != nil {
		return nil, errors.New("unauthorized: invalid user")
	}

	inventories, err := s.inventoryRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return inventories, nil
}

func (s *InventoryService) GetAuthClient() *grpc.AuthClient {
	return s.authClient
}
