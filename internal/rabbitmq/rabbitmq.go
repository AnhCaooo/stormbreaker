package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/AnhCaooo/stormbreaker/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	reconnectDelay    = 5 * time.Second
	reconnectAttempts = 5
)

type RabbitMQ struct {
	config *models.Broker
	ctx    context.Context
	logger *zap.Logger
}

func NewRabbit(ctx context.Context, config *models.Broker, logger *zap.Logger) *RabbitMQ {
	return &RabbitMQ{
		ctx:    ctx,
		config: config,
		logger: logger,
	}
}

// EstablishConnection tries to establish a connection with RabbitMQ server
func (r *RabbitMQ) establishConnection() (*amqp.Connection, error) {
	connection, err := amqp.Dial(r.getURI())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %s", err.Error())
	}
	r.logger.Info("Successfully connected to RabbitMQ")
	return connection, nil
}

// NewProducer retrieves connection client, then opens channel and build producer instance
func (r *RabbitMQ) NewRabbitMQProducer() (*Producer, error) {
	conn, err := r.establishConnection()
	if err != nil {
		return nil, err
	}

	go r.monitorConnection(conn)
	// create a new channel
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel for producer: %s", err.Error())
	}
	return &Producer{
		Channel: ch,
		logger:  r.logger,
		ctx:     r.ctx,
	}, nil
}

// NewConsumer retrieves connection client, then opens channel and build consumer instance
func (r *RabbitMQ) NewRabbitMQConsumer() (*Consumer, error) {
	conn, err := r.establishConnection()
	if err != nil {
		return nil, err
	}
	go r.monitorConnection(conn)
	// create a new channel
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel for consumer: %s", err.Error())
	}

	return &Consumer{
		Channel: ch,
		logger:  r.logger,
	}, nil
}

func (r *RabbitMQ) getURI() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", r.config.Username, r.config.Password, r.config.Host, r.config.Port)
}

// monitorConnection creates a go channel and a goroutine to monitor the connection.
// if connection is lost, then reconnect
func (r RabbitMQ) monitorConnection(conn *amqp.Connection) {
	notifyClose := make(chan *amqp.Error)
	conn.NotifyClose(notifyClose)

	for {
		err := <-notifyClose
		if err != nil {
			r.logger.Warn("Connection closed. Reconnecting...", zap.Error(err))
			var newConn *amqp.Connection
			var reconnectErr error
			for {
				newConn, reconnectErr = r.establishConnection()
				if reconnectErr == nil {
					r.logger.Info("Reconnected to RabbitMQ in goroutine")
					conn = newConn
					notifyClose = make(chan *amqp.Error)
					conn.NotifyClose(notifyClose)
					break
				}
				r.logger.Error(fmt.Sprintf("Reconnection failed. Retrying in %d...", reconnectDelay), zap.Error(reconnectErr))
				time.Sleep(reconnectDelay)
			}
		}
	}
}
