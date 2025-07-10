package service

import (
	"errors"

	"rancher-manager/internal/itemservice/model"
	"rancher-manager/internal/itemservice/repository"
)

type AuthClientInterface interface {
	GetUser(userID uint32) (interface{}, error)
	ValidateToken(token string) (interface{}, error)
}

type ItemService struct {
	itemRepo   *repository.ItemRepository
	authClient AuthClientInterface
}

func NewItemService(itemRepo *repository.ItemRepository, authClient AuthClientInterface) *ItemService {
	return &ItemService{
		itemRepo:   itemRepo,
		authClient: authClient,
	}
}

func (s *ItemService) CreateItem(req *model.CreateItemRequest, userID uint32) (*model.Item, error) {
	// Validate user exists via gRPC
	_, err := s.authClient.GetUser(userID)
	if err != nil {
		return nil, errors.New("unauthorized: invalid user")
	}

	item := &model.Item{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		Stock:       req.Stock,
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}

	if err := s.itemRepo.Create(item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *ItemService) GetItem(id string, userID uint32) (*model.Item, error) {
	// Validate user exists via gRPC
	_, err := s.authClient.GetUser(userID)
	if err != nil {
		return nil, errors.New("unauthorized: invalid user")
	}

	item, err := s.itemRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("item not found")
	}

	return item, nil
}

func (s *ItemService) GetAllItems(userID uint32) ([]*model.Item, error) {
	// Validate user exists via gRPC
	_, err := s.authClient.GetUser(userID)
	if err != nil {
		return nil, errors.New("unauthorized: invalid user")
	}

	items, err := s.itemRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *ItemService) UpdateItem(id string, req *model.UpdateItemRequest, userID uint32) (*model.Item, error) {
	// Validate user exists via gRPC
	_, err := s.authClient.GetUser(userID)
	if err != nil {
		return nil, errors.New("unauthorized: invalid user")
	}

	// Get existing item
	existingItem, err := s.itemRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("item not found")
	}

	// Update fields if provided
	if req.Name != "" {
		existingItem.Name = req.Name
	}
	if req.Description != "" {
		existingItem.Description = req.Description
	}
	if req.Price > 0 {
		existingItem.Price = req.Price
	}
	if req.Category != "" {
		existingItem.Category = req.Category
	}
	if req.Stock >= 0 {
		existingItem.Stock = req.Stock
	}

	existingItem.UpdatedBy = userID

	if err := s.itemRepo.Update(id, existingItem); err != nil {
		return nil, err
	}

	return existingItem, nil
}

func (s *ItemService) DeleteItem(id string, userID uint32) error {
	// Validate user exists via gRPC
	_, err := s.authClient.GetUser(userID)
	if err != nil {
		return errors.New("unauthorized: invalid user")
	}

	// Check if item exists
	_, err = s.itemRepo.GetByID(id)
	if err != nil {
		return errors.New("item not found")
	}

	return s.itemRepo.Delete(id)
}

func (s *ItemService) GetItemsByCategory(category string, userID uint32) ([]*model.Item, error) {
	// Validate user exists via gRPC
	_, err := s.authClient.GetUser(userID)
	if err != nil {
		return nil, errors.New("unauthorized: invalid user")
	}

	items, err := s.itemRepo.GetByCategory(category)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *ItemService) SearchItems(name string, userID uint32) ([]*model.Item, error) {
	// Validate user exists via gRPC
	_, err := s.authClient.GetUser(userID)
	if err != nil {
		return nil, errors.New("unauthorized: invalid user")
	}

	items, err := s.itemRepo.SearchByName(name)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *ItemService) GetAuthClient() AuthClientInterface {
	return s.authClient
}
