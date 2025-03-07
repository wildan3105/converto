package rabbitmq

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"

	"github.com/wildan3105/converto/pkg/logger"
)

type ConnectionManager struct {
	conn            *amqp091.Connection
	channel         *amqp091.Channel
	mu              sync.Mutex
	rabbitURL       string
	notifyClose     chan *amqp091.Error
	ConnectionError chan error
}

const (
	HeartbeatTimeout = 30
	MaxRetries       = 5
	InitialBackoff   = 1 * time.Second
)

var log = logger.GetInstance()

// NewConnectionManager creates a new ConnectionManager and establishes a connection to RabbitMQ with retries
func NewConnectionManager(rabbitURL string) (*ConnectionManager, error) {
	cm := &ConnectionManager{
		rabbitURL:       rabbitURL,
		notifyClose:     make(chan *amqp091.Error),
		ConnectionError: make(chan error),
	}

	if err := cm.connect(); err != nil {
		return nil, err
	}

	go cm.handleReconnect()

	return cm, nil
}

// connect establishes a new RabbitMQ connection and channel with retries.
func (cm *ConnectionManager) connect() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	var conn *amqp091.Connection
	var err error

	for i := 0; i < MaxRetries; i++ {
		conn, err = amqp091.DialConfig(cm.rabbitURL, amqp091.Config{Heartbeat: HeartbeatTimeout})
		if err == nil {
			break
		}
		log.Warn("RabbitMQ connection failed: %v. Retrying in %v...", err, InitialBackoff)
		time.Sleep(time.Duration(i+1) * InitialBackoff)
	}

	if err != nil {
		log.Error("Failed to connect to RabbitMQ after %d retries: %v", MaxRetries, err)
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Error("Failed to open a channel: %v", err)
		return err
	}

	log.Info("Connected to RabbitMQ and opened a channel")
	cm.conn = conn
	cm.channel = ch
	cm.notifyClose = make(chan *amqp091.Error)
	cm.conn.NotifyClose(cm.notifyClose)
	return nil
}

// handleReconnect handles reconnection when the RabbitMQ connection is closed.
func (cm *ConnectionManager) handleReconnect() {
	for {
		err := <-cm.notifyClose
		if err != nil {
			log.Warn("RabbitMQ connection closed, attempting to reconnect: %v", err)
			for {
				if err := cm.connect(); err == nil {
					log.Info("Successfully reconnected to RabbitMQ")
					cm.ConnectionError <- errors.New("connection lost")
					break
				}
				log.Warn("Reconnection attempt failed: %v", err)
				select {
				case <-time.After(InitialBackoff):
				case <-context.Background().Done():
					return
				}
			}
		}
	}
}

func (cm *ConnectionManager) Close() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.channel != nil {
		cm.channel.Close()
	}

	if cm.conn != nil {
		cm.conn.Close()
	}
}

// Ping checks if the connection to RabbitMQ is still open and functioning correctly.
func (cm *ConnectionManager) Ping() error {

	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.conn.IsClosed() {
		return errors.New("RabbitMQ connection is closed")
	}

	err := cm.channel.Qos(1, 0, false)
	if err != nil {
		log.Warn("RabbitMQ connection ping failed: %v", err)
		return err
	}

	return nil
}

// SetupExchangeQueueBinding sets up an exchange, queue, and binding with the routing key.
func (cm *ConnectionManager) SetupExchangeQueueBinding(exchangeName, routingKey, queueName string) error {
	if err := cm.channel.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil); err != nil {
		log.Warn("Failed to declare exchange: %v", err)
		return err
	}

	if _, err := cm.channel.QueueDeclare(queueName, true, false, false, false, nil); err != nil {
		log.Warn("Failed to declare queue: %v", err)
		return err
	}

	if err := cm.channel.QueueBind(queueName, routingKey, exchangeName, false, nil); err != nil {
		log.Warn("Failed to bind queue: %v", err)
		return err
	}

	log.Info("Set up exchange %s, queue %s, with routing key %s", exchangeName, queueName, routingKey)
	return nil
}

// GetChannel safely returns the channel instance.
func (cm *ConnectionManager) GetChannel() *amqp091.Channel {
	return cm.channel
}

// GetConnection safely returns the connection instance.
func (cm *ConnectionManager) GetConnection() *amqp091.Connection {
	return cm.conn
}

// Util functions for testing purpose
// DeleteExchange deletes a RabbitMQ exchange by name
func (cm *ConnectionManager) DeleteExchange(exchangeName string) error {
	err := cm.channel.ExchangeDelete(
		exchangeName, // exchange name
		false,        // ifUnused
		false,        // noWait
	)
	if err != nil {
		log.Warn("Failed to delete exchange %s: %v", exchangeName, err)
		return err
	}
	log.Info("Successfully deleted exchange: %s", exchangeName)
	return nil
}

// DeleteQueue deletes a RabbitMQ queue by name
func (cm *ConnectionManager) DeleteQueue(queueName string) error {
	_, err := cm.channel.QueueDelete(
		queueName, // queue name
		false,     // ifUnused
		false,     // ifEmpty
		false,     // noWait
	)
	if err != nil {
		log.Warn("Failed to delete queue %s: %v", queueName, err)
		return err
	}
	log.Info("Successfully deleted queue: %s", queueName)
	return nil
}

// CheckExchangeExists checks if the exchange exists
func (cm *ConnectionManager) CheckExchangeExists(exchangeName string) (bool, error) {

	err := cm.channel.ExchangeDeclarePassive(
		exchangeName,
		"direct", // same type as when created
		true,     // durable
		false,    // auto-delete
		false,    // internal
		false,    // noWait
		nil,      // arguments
	)

	if err != nil {
		log.Warn("Failed to check exchange %s: %v", exchangeName, err)
		return false, err
	}

	log.Info("Exchange %s exists", exchangeName)
	return true, nil
}

// CheckQueueExists checks if the queue exists
func (cm *ConnectionManager) CheckQueueExists(queueName string) (bool, error) {

	_, err := cm.channel.QueueDeclarePassive(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		log.Warn("Failed to check queue %s: %v", queueName, err)
		return false, err
	}

	log.Info("Queue %s exists", queueName)
	return true, nil
}

// CheckQueueBinding checks if a queue is bound to an exchange using a routing key
func (cm *ConnectionManager) CheckQueueBinding(exchangeName, queueName, routingKey string) (bool, error) {

	err := cm.channel.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		true,         // passive mode (check existence without creating)
		nil,          // arguments
	)

	if err != nil {
		log.Warn("Failed to check binding for exchange %s, queue %s: %v", exchangeName, queueName, err)
		return false, err
	}

	log.Info("Queue %s is bound to exchange %s with routing key %s", queueName, exchangeName, routingKey)
	return true, nil
}
