package utilities

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"tow-management-system-api/repository"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database wraps a Mongo geoplacesClient and logical database handle.
type Database struct {
	client *mongo.Client
	db     *mongo.Database
}

// NewDatabaseConnection creates and verifies a MongoDB connection, returning a Database.
func NewDatabaseConnection() (*Database, error) {

	hostname := os.Getenv("MONGO_CLUSTER_HOSTNAME")
	appName := os.Getenv("APP_NAME")

	var uri string

	if os.Getenv("ENVIRONMENT") == "local" {
		dbUser := os.Getenv("MONGO_DB_USER")
		dbPass := os.Getenv("MONGO_DB_PASSWORD")
		uri = fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority&appName=%s", dbUser, dbPass, hostname, appName)
	} else {
		uri = fmt.Sprintf("mongodb+srv://%s/?authMechanism=MONGODB-AWS&authSource=%%24external&appName=%s", hostname, appName)
	}

	log.Println(uri)

	clientOpts := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("mongo connect error: %w", err)
	}

	db := client.Database("tow-management-system")

	return &Database{
		client: client,
		db:     db,
	}, nil
}

// Close disconnects the geoplacesClient and releases resources.
func (d *Database) Close(ctx context.Context) error {
	if d == nil || d.client == nil {
		return nil
	}
	return d.client.Disconnect(ctx)
}

// DB returns the underlying *mongo.Database (useful for advanced scenarios).
func (d *Database) DB() *mongo.Database {
	return d.db
}

// ----- Repository factories -----
const (
	UserCollection    = "users"
	CompanyCollection = "companies"
	TowCollection     = "tows"
	PriceCollection   = "prices"
)

// CreateUserRepository returns a Mongo-backed user repository.
func (d *Database) CreateUserRepository() *repository.UserMongoRepository {
	coll := UserCollection
	return repository.NewMongoUserRepository(d.db, coll)
}

// CreateCompanyRepository returns a Mongo-backed company repository.
func (d *Database) CreateCompanyRepository() *repository.CompanyMongoRepository {
	coll := CompanyCollection
	return repository.NewMongoCompanyRepository(d.db, coll)
}

// CreateTowRepository returns a tow repository.
func (d *Database) CreateTowRepository() *repository.TowMongoRepository {
	coll := TowCollection
	return repository.NewMongoTowRepository(d.db, coll)
}

// CreatePriceRepository returns a Mongo-backed price repository.
func (d *Database) CreatePriceRepository() *repository.PriceMongoRepository {
	coll := PriceCollection
	return repository.NewMongoPriceRepository(d.db, coll)
}
