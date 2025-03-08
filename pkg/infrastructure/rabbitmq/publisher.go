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

	if err := p.enableConfirmMode(); err != nil {
		log.Error("Failed to enable confirm mode: %v", err)
	}

	return p
}

// enableConfirmMode sets the channel into confirm mode and attaches a listener to handle confirmation
func (p *Publisher) enableConfirmMode() error {
	err := p.channel.Confirm(false)
	if err != nil {
		return err
	}

	go func() {
		for confirmed := range p.channel.NotifyPublish(make(chan amqp091.Confirmation)) {
			if confirmed.Ack {
				log.Info("Message with delivery tag %d confirmed", confirmed.DeliveryTag)
			} else {
				log.Warn("Message with delivery tag %d failed", confirmed.DeliveryTag)
			}
		}
	}()

	log.Info("Enabling confirm mode")

	return nil
}

func (p *Publisher) PublishConversionJob(ctx context.Context, job domain.ConversionJob, exchange, routingKey string) error {
	channel := p.conn.GetChannel()

	jobBytes, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = channel.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		true,  // mandatory
		false, // immediate
		amqp091.Publishing{
			DeliveryMode: amqp091.Persistent,
			ContentType:  "application/json",
			Body:         jobBytes,
			Timestamp:    time.Now(),
		},
	)

	if err != nil {
		return err
	}

	select {
	case confirm := <-channel.NotifyPublish(make(chan amqp091.Confirmation, 1)):
		if !confirm.Ack {
			return fmt.Errorf("failed to confirm publishing")
		}
	case <-ctx.Done():
		return ctx.Err()
	}

	log.Info("Published job %s to exchange %s with routing key %s", job.JobID, exchange, routingKey)
	return nil
}
