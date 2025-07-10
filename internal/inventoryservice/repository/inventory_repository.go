package repository

import (
	"rancher-manager/internal/inventoryservice/model"

	"gorm.io/gorm"
)

type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) Create(inventory *model.Inventory) error {
	return r.db.Create(inventory).Error
}

func (r *InventoryRepository) GetByItemID(itemID string) (*model.Inventory, error) {
	var inventory model.Inventory
	err := r.db.Where("item_id = ?", itemID).First(&inventory).Error
	if err != nil {
		return nil, err
	}
	return &inventory, nil
}

func (r *InventoryRepository) Update(inventory *model.Inventory) error {
	return r.db.Save(inventory).Error
}

func (r *InventoryRepository) Delete(itemID string) error {
	return r.db.Where("item_id = ?", itemID).Delete(&model.Inventory{}).Error
}

func (r *InventoryRepository) GetAll() ([]*model.Inventory, error) {
	var inventories []*model.Inventory
	err := r.db.Find(&inventories).Error
	return inventories, err
}

func (r *InventoryRepository) UpdateStock(itemID string, newStock int, userID uint32) error {
	var inventory model.Inventory
	err := r.db.Where("item_id = ?", itemID).First(&inventory).Error
	if err != nil {
		return err
	}

	inventory.Stock = newStock
	inventory.UpdatedBy = userID

	return r.db.Save(&inventory).Error
}

func (r *InventoryRepository) ExistsByItemID(itemID string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Inventory{}).Where("item_id = ?", itemID).Count(&count).Error
	return count > 0, err
}
