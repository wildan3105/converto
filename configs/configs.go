package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/wildan3105/converto/pkg/logger"
)

// Config holds the environment variables for the app
type Config struct {
	Port                 string `envconfig:"PORT" default:"3000"`
	Environment          string `envconfig:"ENVIRONMENT" default:"local"`
	MongoURI             string `envconfig:"MONGO_URI" required:"true"`
	MongoDbName          string `envconfig:"DB_NAME" required:"true"`
	MongoDbCollection    string `envconfig:"COLLECTION_NAME" required:"true"`
	RabbitMQURI          string `envconfig:"RABBITMQ_URI" required:"true"`
	RabbitMQExchangeName string `envconfig:"RABBITMQ_EXCHANGE_NAME" required:"true"`
	RabbitMQRoutingKey   string `envconfig:"RABBITMQ_ROUTING_KEY" required:"true"`
	RabbitMQQueueName    string `envconfig:"RABBITMQ_QUEUE_NAME" required:"true"`
	BaseDirectory        string `envconfig:"BASE_DIRECTORY" required:"true"`
}

var AppConfig Config

// LoadConfig loads environment variables and handles errors
func LoadConfig() {
	logger := logger.GetInstance()
	err := godotenv.Load() // typically only be used in development environment
	if err != nil {
		logger.Warn("No .env file found or unable to load")
	}

	err = envconfig.Process("", &AppConfig)
	if err != nil {
		log.Fatal("Error loading environment variables: ", err)
	}
}
