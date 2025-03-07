package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository defines database operations for conversions
type Repository interface {
	CreateConversion(ctx context.Context, conversion *Conversion) error
	GetConversion(ctx context.Context, conversionID string) (*Conversion, error)
	UpdateConversion(ctx context.Context, conversionID string, updateData bson.M) error
	ListConversions(ctx context.Context, filter bson.M, limit, offset int64) ([]*Conversion, error)
}

// MongoRepository represents the MongoDB repository
type MongoRepository struct {
	collection *mongo.Collection
}

// NewMongoRepository creates a new instance of MongoRepository
func NewMongoRepository() Repository {
	return &MongoRepository{
		collection: GetConversionCollection(),
	}
}

// CreateConversion inserts a new conversion document
func (r *MongoRepository) CreateConversion(ctx context.Context, conversion *Conversion) error {
	conversion.ID = primitive.NewObjectID()
	conversion.Job.CreatedAt = time.Now()
	conversion.Job.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, conversion)
	return err
}

// GetConversionByID retrieves a conversion document by ID
func (r *MongoRepository) GetConversion(ctx context.Context, conversionID string) (*Conversion, error) {
	var conversion Conversion
	err := r.collection.FindOne(ctx, bson.M{"conversion_id": conversionID}).Decode(&conversion)
	if err != nil {
		return nil, err
	}

	return &conversion, nil
}

// UpdateConversion updates a conversion document by ID
func (r *MongoRepository) UpdateConversion(ctx context.Context, conversionID string, updateData bson.M) error {
	filter := bson.M{"conversion_id": conversionID}
	update := bson.M{"$set": updateData, "$currentDate": bson.M{"job.updated_at": true}}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// ListConversions retrieves a list of conversion documents
func (r *MongoRepository) ListConversions(ctx context.Context, filter bson.M, limit, offset int64) ([]*Conversion, error) {
	var conversions []*Conversion

	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSkip(offset)
	findOptions.SetSort(bson.D{{"job.created_at", -1}})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var conversion Conversion
		if err := cursor.Decode(&conversion); err != nil {
			return nil, err
		}

		conversions = append(conversions, &conversion)
	}

	return conversions, nil
}
