package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/wildan3105/converto/pkg/api/schema"
)

// Consumer is responsible for consuming messages from RabbitMQ
type Consumer struct {
	conn *ConnectionManager
}

// NewConsumer creates a new Consumer
func NewConsumer(cm *ConnectionManager) *Consumer {
	return &Consumer{
		conn: cm,
	}
}

func (c *Consumer) Consume(ctx context.Context, queueName string) (<-chan schema.ConversionEvent, error) {
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
		return nil, fmt.Errorf("failed to consume messages: %w", err)
	}

	jobChan := make(chan schema.ConversionEvent)

	go func() {
		defer close(jobChan)
		for {
			select {
			case <-ctx.Done():
				log.Info("Context cancelled, stopping message consumption...")
				return
			case msg, ok := <-msgs:
				if !ok {
					log.Warn("Message channel closed by RabbitMQ")
					return
				}

				job := schema.ConversionEvent{}
				if err := json.Unmarshal(msg.Body, &job); err != nil {
					log.Warn("Failed to unmarshal job: %v", err)
					_ = msg.Nack(false, true)
					continue
				}

				jobChan <- job

				if err := msg.Ack(false); err != nil {
					log.Warn("Failed to ack message: %v", err)
				}
				fmt.Printf("Processed job %v", job)
			}
		}
	}()

	return jobChan, nil
}
