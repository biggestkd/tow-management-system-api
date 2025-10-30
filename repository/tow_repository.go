package repository

import (
	"context"
	"fmt"
	"tow-management-system-api/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// TowMongoRepository handles MongoDB operations for the Tow model.
type TowMongoRepository struct {
	collection *mongo.Collection
}

// NewMongoTowRepository creates a new TowMongoRepository instance.
func NewMongoTowRepository(db *mongo.Database, collectionName string) *TowMongoRepository {
	return &TowMongoRepository{
		collection: db.Collection(collectionName),
	}
}

// Create inserts a new tow document into MongoDB.
func (r *TowMongoRepository) Create(ctx context.Context, tow *model.Tow) error {
	_, err := r.collection.InsertOne(ctx, tow)
	if err != nil {
		return fmt.Errorf("failed to create tow: %w", err)
	}
	return nil
}

// Find retrieves tows matching the provided filter struct.
func (r *TowMongoRepository) Find(ctx context.Context, filterModel *model.Tow) ([]*model.Tow, error) {
	bsonBytes, err := bson.Marshal(filterModel)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tow filter: %w", err)
	}

	var filter bson.M
	if err := bson.Unmarshal(bsonBytes, &filter); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tow filter: %w", err)
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find tows: %w", err)
	}
	defer cursor.Close(ctx)

	var results []*model.Tow
	for cursor.Next(ctx) {
		var t model.Tow
		if err := cursor.Decode(&t); err != nil {
			return nil, fmt.Errorf("failed to decode tow document: %w", err)
		}
		results = append(results, &t)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor iteration error: %w", err)
	}

	return results, nil
}

// Update modifies a tow document by ID.
func (r *TowMongoRepository) Update(ctx context.Context, id string, updateData *model.Tow) error {
	bsonBytes, err := bson.Marshal(updateData)
	if err != nil {
		return fmt.Errorf("failed to marshal tow update data: %w", err)
	}

	var updateFields bson.M
	if err := bson.Unmarshal(bsonBytes, &updateFields); err != nil {
		return fmt.Errorf("failed to unmarshal tow update fields: %w", err)
	}

	// Never allow updating the _id via $set
	delete(updateFields, "_id")

	update := bson.M{"$set": updateFields}
	filter := bson.M{"_id": id}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update tow: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("tow with id %s not found", id)
	}

	return nil
}

// Delete removes a tow document by ID.
func (r *TowMongoRepository) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}

	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete tow: %w", err)
	}

	return nil
}
