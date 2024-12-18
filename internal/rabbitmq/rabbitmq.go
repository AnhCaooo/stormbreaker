package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/AnhCaooo/stormbreaker/internal/constants"
	"github.com/AnhCaooo/stormbreaker/internal/db"
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
	mongo  *db.Mongo
}

func NewRabbit(ctx context.Context, config *models.Broker, logger *zap.Logger, mongo *db.Mongo) *RabbitMQ {
	return &RabbitMQ{
		ctx:    ctx,
		config: config,
		logger: logger,
		mongo:  mongo,
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

/*
-------------------------------- PRODUCER METHODS --------------------------------
*/

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

/*
-------------------------------- CONSUMER METHODS --------------------------------
*/

// NewConsumer retrieves connection client, then opens channel and build consumer instance
func (r *RabbitMQ) newRabbitMQConsumer() (*Consumer, error) {
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
		mongo:   r.mongo,
	}, nil
}

// StartRabbitMQConsumer starts the RabbitMQ consumer based on given queue name
func (r *RabbitMQ) StartRabbitMQConsumer(queueName string) {
	messageConsumer, err := r.newRabbitMQConsumer()
	if err != nil {
		r.logger.Fatal(constants.Server, zap.Error(err))
	}
	defer messageConsumer.Channel.Close()
	messageConsumer.Listen(queueName)
}