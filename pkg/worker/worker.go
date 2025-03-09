package worker

import (
	"context"

	"github.com/wildan3105/converto/pkg/domain"
	"github.com/wildan3105/converto/pkg/infrastructure/filestorage"
	"github.com/wildan3105/converto/pkg/infrastructure/rabbitmq"
	"github.com/wildan3105/converto/pkg/logger"
	"github.com/wildan3105/converto/pkg/repository"
)

var log = logger.GetInstance()

// Worker is the core struct for managing job consumption and processing
type Worker struct {
	consumer *rabbitmq.Consumer
	repo     repository.ConversionRepository
	storage  filestorage.FileStorage
}

// NewWorker creates a new Worker instance
func NewWorker(consumer *rabbitmq.Consumer, repo repository.ConversionRepository, storage filestorage.FileStorage) *Worker {
	return &Worker{
		consumer: consumer,
		repo:     repo,
		storage:  storage,
	}
}

// Start begins consuming messages and processing conversion jobs
func (w *Worker) Start(ctx context.Context, queueName string) error {
	jobChan, err := w.consumer.Consume(ctx, queueName)
	if err != nil {
		log.Error("Failed to start consuming messages: %v", err)
		return err
	}

	for {
		select {
		case <-ctx.Done():
			log.Info("Worker context cancelled, stopping job processing...")
			return nil
		case job, ok := <-jobChan:
			if !ok {
				log.Info("Job channel closed, exiting worker...")
				return nil
			}
			go func(job domain.ConversionJob) {
				if err := w.Handle(ctx, job); err != nil {
					log.Error("Job processing failed: %v", err)
				}
			}(job)
		}
	}
}
