package repository

import (
	"context"
	"fmt"
	"tow-management-system-api/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// PriceMongoRepository handles MongoDB operations for the Price model.
type PriceMongoRepository struct {
	collection *mongo.Collection
}

// NewMongoPriceRepository creates a new PriceMongoRepository instance.
func NewMongoPriceRepository(db *mongo.Database, collectionName string) *PriceMongoRepository {
	return &PriceMongoRepository{
		collection: db.Collection(collectionName),
	}
}

// Create inserts a new price document into MongoDB.
func (r *PriceMongoRepository) Create(ctx context.Context, price *model.Price) error {
	_, err := r.collection.InsertOne(ctx, price)
	if err != nil {
		return fmt.Errorf("failed to create price: %w", err)
	}
	return nil
}

// Find retrieves prices matching the provided filter struct.
func (r *PriceMongoRepository) Find(ctx context.Context, filterModel *model.Price) ([]*model.Price, error) {
	bsonBytes, err := bson.Marshal(filterModel)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal price filter: %w", err)
	}

	var filter bson.M
	if err := bson.Unmarshal(bsonBytes, &filter); err != nil {
		return nil, fmt.Errorf("failed to unmarshal price filter: %w", err)
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find prices: %w", err)
	}
	defer cursor.Close(ctx)

	var results []*model.Price
	for cursor.Next(ctx) {
		var p model.Price
		if err := cursor.Decode(&p); err != nil {
			return nil, fmt.Errorf("failed to decode price document: %w", err)
		}
		results = append(results, &p)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor iteration error: %w", err)
	}

	return results, nil
}

// Update modifies a price document by ID.
func (r *PriceMongoRepository) Update(ctx context.Context, id string, updateData *model.Price) error {
	bsonBytes, err := bson.Marshal(updateData)
	if err != nil {
		return fmt.Errorf("failed to marshal price update data: %w", err)
	}

	var updateFields bson.M
	if err := bson.Unmarshal(bsonBytes, &updateFields); err != nil {
		return fmt.Errorf("failed to unmarshal price update fields: %w", err)
	}

	// Never allow updating the _id via $set
	delete(updateFields, "_id")

	update := bson.M{"$set": updateFields}
	filter := bson.M{"_id": id}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update price: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("price with id %s not found", id)
	}

	return nil
}

// Delete removes a price document by ID.
func (r *PriceMongoRepository) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}

	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete price: %w", err)
	}

	return nil
}
