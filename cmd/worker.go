package cmd

import (
	externalLog "log"

	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	config "github.com/wildan3105/converto/configs"
	"github.com/wildan3105/converto/pkg/infrastructure/filestorage"
	"github.com/wildan3105/converto/pkg/infrastructure/mongodb"
	"github.com/wildan3105/converto/pkg/infrastructure/rabbitmq"
	"github.com/wildan3105/converto/pkg/repository"
	rabbitMQWorker "github.com/wildan3105/converto/pkg/worker"
)

// WorkerCMd is the command to run the worker (background job processor)
var WorkerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Run the worker",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Starting worker process...")

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			sig := <-sigChan
			log.Info("Signal received: %s. Shutting down worker...", sig)
			cancel()
		}()

		mongoClient, err := mongodb.Connect(config.AppConfig.MongoURI)
		if err != nil {
			externalLog.Fatal("Failed to connect to MongoDB: ", err)
		}

		connManager, err := rabbitmq.NewConnectionManager(config.AppConfig.RabbitMQURI)
		if err != nil {
			externalLog.Fatalf("Failed to initialize RabbitMQ: %v", err)
		}

		consumer := rabbitmq.NewConsumer(connManager)
		conversionRepo := repository.NewMongoRepository(mongoClient, config.AppConfig.MongoDbName)
		storage := filestorage.NewLocalFileStorage(config.AppConfig.BaseDirectory)

		worker := rabbitMQWorker.NewWorker(consumer, conversionRepo, storage)

		workerErr := worker.Start(ctx, "conversion_queue")

		if workerErr != nil {
			externalLog.Fatalf("Failed when consuming messages: %v", workerErr)
		}

		<-ctx.Done()

		log.Info("Worker process exiting gracefully...")
	},
}
