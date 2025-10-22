package repository

import (
	"context"
	"fmt"
	"log"
	"tow-management-system-api/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserMongoRepository handles MongoDB operations for the User model.
type UserMongoRepository struct {
	collection *mongo.Collection
}

// NewMongoUserRepository creates a new UserMongoRepository instance.
func NewMongoUserRepository(db *mongo.Database, collectionName string) *UserMongoRepository {
	return &UserMongoRepository{
		collection: db.Collection(collectionName),
	}
}

// Create inserts a new user document into MongoDB.
func (r *UserMongoRepository) Create(ctx context.Context, user *model.User) error {
	log.Println("Running UserMongoRepository Create")
	log.Printf("%v", *user)

	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// Find retrieves users matching the provided filter struct.
func (r *UserMongoRepository) Find(ctx context.Context, filterModel *model.User) ([]*model.User, error) {
	bsonBytes, err := bson.Marshal(filterModel)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user filter: %w", err)
	}

	var filter bson.M
	if err := bson.Unmarshal(bsonBytes, &filter); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user filter: %w", err)
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}
	defer cursor.Close(ctx)

	var results []*model.User
	for cursor.Next(ctx) {
		var u model.User
		if err := cursor.Decode(&u); err != nil {
			return nil, fmt.Errorf("failed to decode user document: %w", err)
		}
		results = append(results, &u)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor iteration error: %w", err)
	}

	return results, nil
}

// Update modifies a user document by ID.
func (r *UserMongoRepository) Update(ctx context.Context, id string, updateData *model.User) error {

	bsonBytes, err := bson.Marshal(updateData)
	if err != nil {
		return fmt.Errorf("failed to marshal user update data: %w", err)
	}

	var updateFields bson.M
	if err := bson.Unmarshal(bsonBytes, &updateFields); err != nil {
		return fmt.Errorf("failed to unmarshal user update fields: %w", err)
	}

	// Remove immutable fields
	delete(updateFields, "createdDate")
	delete(updateFields, "_id")

	update := bson.M{
		"$set": updateFields,
	}

	filter := bson.M{"_id": id}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("user with id %s not found", id)
	}

	return nil
}

// Delete removes a user document by ID.
func (r *UserMongoRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
