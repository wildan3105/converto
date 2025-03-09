package service

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/wildan3105/converto/pkg/infrastructure/mongodb"
	"github.com/wildan3105/converto/pkg/infrastructure/rabbitmq"
)

// HealthService provides methods for health checks
type HealthService struct {
	MongoClient    *mongo.Client
	RabbitMQClient *rabbitmq.ConnectionManager
}

// NewHealthService returns a new HealthService
func NewHealthService(mongoClient *mongo.Client, rabbitMQClient *rabbitmq.ConnectionManager) *HealthService {
	return &HealthService{
		MongoClient:    mongoClient,
		RabbitMQClient: rabbitMQClient,
	}
}

// CheckDependencies pings the dependencies (MongoDB and RabbitMQ)
func (hs *HealthService) CheckDependencies() error {
	if err := mongodb.Ping(hs.MongoClient); err != nil {
		return fmt.Errorf("MongoDB ping failed: %w", err)
	}

	if err := hs.RabbitMQClient.Ping(); err != nil {
		return fmt.Errorf("RabbitMQ ping failed: %w", err)
	}

	return nil
}
