package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/wildan3105/converto/pkg/api/schema"
)

// Publisher is responsible for publishing messages to RabbitMQ
type Publisher struct {
	connManager *ConnectionManager
	channel     *amqp091.Channel
}

// NewPublisher creates a new Publisher
func NewPublisher(cm *ConnectionManager) *Publisher {
	p := &Publisher{
		connManager: cm,
		channel:     cm.channel,
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

func (p *Publisher) PublishConversionJob(ctx context.Context, conversionEvent schema.ConversionEvent, exchange, routingKey string) error {
	jobBytes, err := json.Marshal(conversionEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	if p.channel.IsClosed() {
		log.Warn("RabbitMQ channel is closed, attempting to re-open it.")
		if err := p.reOpenChannel(); err != nil {
			log.Error("Failed to re-open RabbitMQ channel: %v", err)
			return amqp091.ErrClosed
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	mandatory := true
	returnCh := make(chan amqp091.Return)
	p.channel.NotifyReturn(returnCh)

	go func() {
		for returned := range returnCh {
			log.Warn("Message returned: %s", returned.ReplyText)
		}
	}()

	err = p.channel.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		mandatory,
		false, // immediate
		amqp091.Publishing{
			DeliveryMode: amqp091.Persistent,
			ContentType:  "application/json",
			Body:         jobBytes,
			Timestamp:    time.Now(),
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	select {
	case <-ctx.Done():
		log.Warn("Publishing context timed out or cancelled")
		return ctx.Err()
	default:
		log.Info("Published job %s to exchange %s with routing key %s", conversionEvent.JobID, exchange, routingKey)
		return nil
	}
}

// reOpenChannel attempts to re-open the RabbitMQ channel
func (p *Publisher) reOpenChannel() error {
	var err error
	p.channel, err = p.connManager.conn.Channel()
	if err != nil {
		return err
	}

	log.Info("Attempting to re-open channel")
	return p.enableConfirmMode()
}
