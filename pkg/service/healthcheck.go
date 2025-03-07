package service

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/wildan3105/converto/pkg/infrastructure/mongodb"
)

// HealthService provides methods for health checks
type HealthService struct {
	MongoClient *mongo.Client
}

// NewHealthService returns a new HealthService
func NewHealthService(mongoClient *mongo.Client) *HealthService {
	return &HealthService{
		MongoClient: mongoClient,
	}
}

// CheckDependencies pings the dependencies (MongoDB and RabbitMQ)
func (hs *HealthService) CheckDependencies() error {
	if err := mongodb.Ping(hs.MongoClient); err != nil {
		return fmt.Errorf("MongoDB ping failed: %w", err)
	}

	return nil
}
