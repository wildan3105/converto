package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/wildan3105/converto/pkg/domain"
	"github.com/wildan3105/converto/pkg/infrastructure/circuitbreaker"
	"github.com/wildan3105/converto/pkg/infrastructure/mongodb"

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
	collection     *mongo.Collection
	circuitbreaker *circuitbreaker.CircuitBreaker
}

// NewMongoRepository creates a new instance of MongoRepository
func NewMongoRepository(mongoClient *mongo.Client, dbName string) *MongoConversionRepository {
	cb := circuitbreaker.NewCircuitBreaker(3, 10*time.Second) // 3 failures, 10second cooldown
	return &MongoConversionRepository{
		collection:     mongoClient.Database(dbName).Collection("conversions"),
		circuitbreaker: cb,
	}
}

// CreateConversion inserts a new conversion document
func (r *MongoConversionRepository) CreateConversion(ctx context.Context, conversion *domain.Conversion) (string, error) {
	ctx, cancel := mongodb.WithTimeout(ctx)
	defer cancel()

	conversion.ID = primitive.NewObjectID().Hex()
	conversion.Job.CreatedAt = time.Now()
	conversion.Job.UpdatedAt = time.Now()

	err := r.circuitbreaker.Execute(func() error {
		_, err := r.collection.InsertOne(ctx, conversion)
		if err != nil {
			if mongo.IsTimeout(err) || errors.Is(err, context.DeadlineExceeded) {
				fmt.Println("Error: database operation timed out")
				return errors.New("database operation timed out")
			}
			return err
		}
		return nil
	})

	if err != nil {
		fmt.Printf("CreateConversion failed: %v\n", err)
		return "", err
	}

	return conversion.ID, nil
}

// GetConversionByID retrieves a conversion document by ID
func (r *MongoConversionRepository) GetConversionByID(ctx context.Context, conversionID string) (*domain.Conversion, error) {
	ctx, cancel := mongodb.WithTimeout(ctx)
	defer cancel()

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
	ctx, cancel := mongodb.WithTimeout(ctx)
	defer cancel()

	filter := bson.M{"_id": conversionID}
	update := bson.M{
		"$set":         updateData,
		"$currentDate": bson.M{"job.updatedAt": true},
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
	ctx, cancel := mongodb.WithTimeout(ctx)
	defer cancel()

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
