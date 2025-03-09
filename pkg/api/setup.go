package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	config "github.com/wildan3105/converto/configs"
	"github.com/wildan3105/converto/pkg/handler"
	"github.com/wildan3105/converto/pkg/infrastructure/filestorage"
	"github.com/wildan3105/converto/pkg/infrastructure/mongodb"
	"github.com/wildan3105/converto/pkg/infrastructure/rabbitmq"
	"github.com/wildan3105/converto/pkg/repository"
	"github.com/wildan3105/converto/pkg/service"
)

func Setup() *fiber.App {
	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
		Next: func(c *fiber.Ctx) bool {
			return c.Path() == "/api/health"
		},
	}))

	mongoClient, err := mongodb.Connect(config.AppConfig.MongoURI)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
	}

	connManager, err := rabbitmq.NewConnectionManager(config.AppConfig.RabbitMQURI)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}

	publisher := rabbitmq.NewPublisher(connManager)

	conversionRepo := repository.NewMongoRepository(mongoClient, config.AppConfig.MongoDbName)
	storage := filestorage.NewLocalFileStorage(config.AppConfig.BaseDirectory)

	conversionService := service.NewConversionService(conversionRepo, publisher, storage)
	healthService := service.NewHealthService(mongoClient)

	conversionHandler := handler.NewConversionHandler(conversionService)
	healthHandler := handler.NewHealthHandler(healthService)

	api := app.Group("/api")

	api.Get("/health", healthHandler.Check)

	v1 := api.Group("/v1")
	v1.Post("/conversions", conversionHandler.CreateConversion)
	v1.Get("/conversions", conversionHandler.GetConversions)
	v1.Get("/conversions/:id", conversionHandler.GetConversionByID)
	// v1.Get("/files/original/:conversionId", conversionHandler.GetOriginalFileByConversionId)
	// v1.Get("/files/converted/:conversionId", conversionHandler.GetConvertedFileByConversionId)

	return app
}
