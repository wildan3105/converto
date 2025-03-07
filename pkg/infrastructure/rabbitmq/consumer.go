package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rabbitmq/amqp091-go"

	"github.com/wildan3105/converto/pkg/domain"
)

// Consumer is responsible for consuming messages from RabbitMQ
type Consumer struct {
	conn    *ConnectionManager
	channel *amqp091.Channel
}

// NewConsumer creates a new Consumer
func NewConsumer(cm *ConnectionManager) *Consumer {
	c := &Consumer{
		conn:    cm,
		channel: cm.GetChannel(),
	}

	return c
}

// Consume consumes messages from the given queue and processes them with the provided handler
func (c *Consumer) Consume(ctx context.Context, queueName string, handler func(domain.ConversionJob) error) error {
	channel := c.conn.GetChannel()

	msgs, err := channel.Consume(
		queueName,
		"",
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	if err != nil {
		return fmt.Errorf("failed to consume messages: %w", err)
	}

	go func() {
		for msg := range msgs {
			job := domain.ConversionJob{}
			if err := json.Unmarshal(msg.Body, &job); err != nil {
				log.Warn("Failed to unmarshal job: %v", err)
				if err := msg.Nack(false, true); err != nil {
					log.Warn("Failed to nack message: %v", err)
				}
				continue
			}

			if err := handler(job); err != nil {
				log.Warn("Failed to process job: %v", err)
				if err := msg.Nack(false, true); err != nil {
					log.Warn("Failed to nack message: %v", err)
				}
				continue
			}

			if err := msg.Ack(false); err != nil {
				log.Warn("Failed to ack message: %v", err)
			}
			log.Info("Processed job %s", job.JobID)
		}
	}()

	log.Info("Started consuming messages from queue %s", queueName)
	<-ctx.Done()
	return nil
}
