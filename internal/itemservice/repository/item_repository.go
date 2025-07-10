package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"rancher-manager/internal/itemservice/model"
)

type ItemRepository struct {
	collection *mongo.Collection
}

func NewItemRepository(db *mongo.Database) *ItemRepository {
	return &ItemRepository{
		collection: db.Collection("items"),
	}
}

func (r *ItemRepository) Create(item *model.Item) error {
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(context.Background(), item)
	if err != nil {
		return err
	}

	item.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ItemRepository) GetByID(id string) (*model.Item, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var item model.Item
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *ItemRepository) GetAll() ([]*model.Item, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var items []*model.Item
	if err = cursor.All(context.Background(), &items); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *ItemRepository) Update(id string, item *model.Item) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	item.UpdatedAt = time.Now()
	item.ID = objectID

	_, err = r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objectID},
		bson.M{"$set": item},
	)
	return err
}

func (r *ItemRepository) Delete(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	return err
}

func (r *ItemRepository) GetByCategory(category string) ([]*model.Item, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{"category": category})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var items []*model.Item
	if err = cursor.All(context.Background(), &items); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *ItemRepository) SearchByName(name string) ([]*model.Item, error) {
	filter := bson.M{"name": bson.M{"$regex": name, "$options": "i"}}

	cursor, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var items []*model.Item
	if err = cursor.All(context.Background(), &items); err != nil {
		return nil, err
	}

	return items, nil
}
