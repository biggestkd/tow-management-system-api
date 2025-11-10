package repository

import (
	"context"
	"fmt"
	"tow-management-system-api/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CompanyMongoRepository handles MongoDB operations for the Company model.
type CompanyMongoRepository struct {
	collection *mongo.Collection
}

// NewMongoCompanyRepository creates a new CompanyMongoRepository instance.
func NewMongoCompanyRepository(db *mongo.Database, collectionName string) *CompanyMongoRepository {
	return &CompanyMongoRepository{
		collection: db.Collection(collectionName),
	}
}

// Create inserts a new company document into MongoDB.
func (r *CompanyMongoRepository) Create(ctx context.Context, company *model.Company) error {
	_, err := r.collection.InsertOne(ctx, company)
	if err != nil {
		return fmt.Errorf("failed to create company: %w", err)
	}
	return nil
}

// Find retrieves companies matching the provided filter struct.
func (r *CompanyMongoRepository) Find(ctx context.Context, filterModel *model.Company) ([]*model.Company, error) {
	bsonBytes, err := bson.Marshal(filterModel)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal company filter: %w", err)
	}

	var filter bson.M
	if err := bson.Unmarshal(bsonBytes, &filter); err != nil {
		return nil, fmt.Errorf("failed to unmarshal company filter: %w", err)
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find companies: %w", err)
	}
	defer cursor.Close(ctx)

	var results []*model.Company
	for cursor.Next(ctx) {
		var c model.Company
		if err := cursor.Decode(&c); err != nil {
			return nil, fmt.Errorf("failed to decode company document: %w", err)
		}
		results = append(results, &c)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor iteration error: %w", err)
	}

	return results, nil
}

// Update modifies a company document by ID.
func (r *CompanyMongoRepository) Update(ctx context.Context, id string, updateData *model.Company) error {

	bsonBytes, err := bson.Marshal(updateData)

	if err != nil {
		return fmt.Errorf("failed to marshal company update data: %w", err)
	}

	var updateFields bson.M

	if err := bson.Unmarshal(bsonBytes, &updateFields); err != nil {
		return fmt.Errorf("failed to unmarshal company update fields: %w", err)
	}

	// Never allow updating the _id via $set
	delete(updateFields, "_id")

	update := bson.M{
		"$set": updateFields,
	}

	filter := bson.M{"_id": id}

	result, err := r.collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return fmt.Errorf("failed to update company: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("company with id %s not found", id)
	}

	return nil

}

// Delete removes a company document by ID.
func (r *CompanyMongoRepository) Delete(ctx context.Context, id string) error {

	filter := bson.M{"_id": id}

	_, err := r.collection.DeleteOne(ctx, filter)

	if err != nil {
		return fmt.Errorf("failed to delete company: %w", err)
	}

	return nil
}
