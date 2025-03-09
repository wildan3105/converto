package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/wildan3105/converto/pkg/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConversionRepository defines database operations for conversions
type ConversionRepository interface {
	CreateConversion(ctx context.Context, conversion *domain.Conversion) (string, error)
	GetConversionByID(ctx context.Context, conversionID string) (*domain.Conversion, error)
	UpdateConversion(ctx context.Context, conversionID string, updateData bson.M) error
	ListConversions(ctx context.Context, status string, limit, offset int) ([]*domain.Conversion, error)
}

// MongoConversionRepository represents the MongoDB repository
type MongoConversionRepository struct {
	collection *mongo.Collection
}

// NewMongoRepository creates a new instance of MongoRepository
func NewMongoRepository(mongoClient *mongo.Client, dbName string) *MongoConversionRepository {
	return &MongoConversionRepository{
		collection: mongoClient.Database(dbName).Collection("conversions"),
	}
}

// CreateConversion inserts a new conversion document
func (r *MongoConversionRepository) CreateConversion(ctx context.Context, conversion *domain.Conversion) (string, error) {
	conversion.ID = primitive.NewObjectID().Hex()
	conversion.Job.CreatedAt = time.Now()
	conversion.Job.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, conversion)
	if err != nil {
		return "", err
	}

	return conversion.ID, nil
}

// GetConversionByID retrieves a conversion document by ID
func (r *MongoConversionRepository) GetConversionByID(ctx context.Context, conversionID string) (*domain.Conversion, error) {
	var conversion domain.Conversion
	err := r.collection.FindOne(ctx, bson.M{"_id": conversionID}).Decode(&conversion)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &conversion, nil
}

// UpdateConversion updates a conversion document by ID
func (r *MongoConversionRepository) UpdateConversion(ctx context.Context, conversionID string, updateData bson.M) error {
	filter := bson.M{"_id": conversionID} // Use the correct filter field
	update := bson.M{
		"$set":         updateData,
		"$currentDate": bson.M{"job.updated_at": true}, // Automatically set to the current date
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		fmt.Println("No document matched the provided conversion ID.")
		return errors.New("conversion not found")
	}

	if res.ModifiedCount > 0 {
		fmt.Println("Successfully updated the conversion document!")
	} else {
		fmt.Println("Document found but no update was necessary.")
	}

	return nil
}

// ListConversions retrieves a list of conversion documents with optional status filtering
func (r *MongoConversionRepository) ListConversions(ctx context.Context, status string, limit, offset int) ([]*domain.Conversion, error) {
	var conversions []*domain.Conversion

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	filter := bson.M{}
	if status != "" {
		filter["conversion.status"] = status
	}

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var conversion domain.Conversion
		if err := cursor.Decode(&conversion); err != nil {
			return nil, err
		}
		conversions = append(conversions, &conversion)
	}

	return conversions, nil
}
