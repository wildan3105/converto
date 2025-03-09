package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	config "github.com/wildan3105/converto/configs"
	"github.com/wildan3105/converto/pkg/logger"
)

var MongoClient *mongo.Client

const DefaultTimeout = 10 * time.Second

var log = logger.GetInstance()

// Connect initializes a MongoDB connection and returns the client
// It includes reconnection logic that retries up to 5 times with a 5-second interval
func Connect(uri string) (*mongo.Client, error) {
	var client *mongo.Client
	var err error

	clientOptions := options.Client().ApplyURI(config.AppConfig.MongoURI)

	for attempt := 1; attempt <= 5; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		client, err = mongo.Connect(ctx, clientOptions)
		if err == nil {
			err = client.Ping(ctx, nil)
		}

		cancel()

		if err == nil {
			log.Info("Connected to MongoDB")
			MongoClient = client
			return client, nil
		}

		log.Warn("Failed to connect to MongoDB (attempt %d): %v\n", attempt, err)

		if attempt < 5 {
			log.Info("Retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
		}
	}

	return nil, fmt.Errorf("failed to connect to MongoDB after 5 attempts: %w", err)
}

// Ping checks for the connection status
func Ping(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return client.Ping(ctx, nil)
}

// Disconnect disconnects from MongoDB
func Disconnect(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect MongoDB: %w", err)
	}

	log.Info("Disconnected from MongoDB")
	return nil
}

// GetCollection returns a MongoDB collection
func GetCollection(client *mongo.Client, databaseName, collectionName string) *mongo.Collection {
	return client.Database(databaseName).Collection(collectionName)
}

// WithTimeout applies a generic timeout to a context for MongoDB operations
func WithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, DefaultTimeout)
}
