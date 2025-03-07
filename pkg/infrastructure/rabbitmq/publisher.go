package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/wildan3105/converto/pkg/domain"
)

// Publisher is responsible for publishing messages to RabbitMQ
type Publisher struct {
	conn    *ConnectionManager
	channel *amqp091.Channel
}

// NewPublisher creates a new Publisher
func NewPublisher(cm *ConnectionManager) *Publisher {
	p := &Publisher{
		conn:    cm,
		channel: cm.GetChannel(),
	}

	return p
}

// PublishConversionJob publishes a conversion job to RabbitMQ
func (p *Publisher) PublishConversionJob(ctx context.Context, job domain.ConversionJob, exchange, routingKey string) error {
	channel := p.conn.GetChannel()

	jobBytes, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	err = channel.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        jobBytes,
			Timestamp:   time.Now(),
		},
	)

	if err != nil {
		log.Warn("Failed to publish job: %v", err)
		return err
	}

	log.Info("Published job %s to exchange %s with routing key %s", job.JobID, exchange, routingKey)
	return nil
}
