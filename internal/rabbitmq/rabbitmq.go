// AnhCao 2024
//
// Package rabbitmq provides a RabbitMQ client implementation for establishing connections,
// creating producers and consumers, and managing message queues. It includes functionality
// for monitoring and reconnecting to RabbitMQ servers, as well as handling message
// production and consumption with multiple workers.
//
// Types:
//
// - RabbitMQ: Main struct for managing RabbitMQ connections and channels
//
// - Producer: Struct for managing message production
//
// - Consumer: Struct for managing message consumption
package rabbitmq

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"

	"github.com/AnhCaooo/stormbreaker/internal/db"
	"github.com/AnhCaooo/stormbreaker/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	exchange_type = "topic" // exchange type for RabbitMQ
)

// RabbitMQ represents a RabbitMQ broker instance with its configuration,
// connection, channels, context, logger, and MongoDB instance.
type RabbitMQ struct {
	// Configuration settings for the RabbitMQ broker.
	config *models.Broker
	// The AMQP connection to the RabbitMQ server.
	connection *amqp.Connection
	// A slice of AMQP channels for communication with RabbitMQ.
	channels []*amqp.Channel
	// The context for managing the lifecycle of the RabbitMQ instance.
	ctx context.Context
	// The logger for logging RabbitMQ-related activities
	logger *zap.Logger
	// The MongoDB instance for database operations.
	mongo *db.Mongo
	//  A channel to send errors encountered during the consumer setup and operation.
	errChan chan<- error
	// A channel to signal the consumer to stop listening for messages.
	stopChan <-chan struct{}
	// A pointer to a sync.WaitGroup to signal when the consumer has finished.
	wg *sync.WaitGroup
}

// NewRabbit creates a new instance of RabbitMQ with the provided context, configuration, logger, and MongoDB client.
// It initializes the RabbitMQ struct with the given parameters.
func NewRabbit(ctx context.Context, config *models.Broker, logger *zap.Logger, mongo *db.Mongo) *RabbitMQ {
	return &RabbitMQ{
		ctx:    ctx,
		config: config,
		logger: logger,
		mongo:  mongo,
	}
}

// EstablishConnection establishes a connection with RabbitMQ server.
// Returns an error if the connection fails.
func (r *RabbitMQ) EstablishConnection() (err error) {
	r.connection, err = amqp.Dial(r.getURI())
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %s", err.Error())
	}
	r.logger.Info("successfully connected to RabbitMQ")
	return nil
}

func (r *RabbitMQ) getURI() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", r.config.Username, r.config.Password, r.config.Host, r.config.Port)
}

// CloseConnection closes first all channels then the connection with RabbitMQ server.
func (r *RabbitMQ) CloseConnection() {
	// Close all channels
	for _, channel := range r.channels {
		if err := channel.Close(); err != nil {
			errMsg := fmt.Errorf("failed to close channel: %s", err.Error())
			r.logger.Fatal(errMsg.Error())
		}
	}
	r.logger.Info("closed RabbitMQ channels")
	// Close the connection
	if r.connection != nil {
		if err := r.connection.Close(); err != nil {
			errMsg := fmt.Errorf("failed to close RabbitMQ connection: %s", err.Error())
			r.logger.Fatal(errMsg.Error())
		}
	}
	r.logger.Info("closed RabbitMQ connection")
}

/*
-------------------------------- PRODUCER METHODS --------------------------------
*/

// NewProducer retrieves connection client, then opens channel and build producer instance
// func (r *RabbitMQ) newProducer() (*Producer, error) {
// 	// create a new channel
// 	ch, err := r.connection.Channel()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to open a channel for producer: %s", err.Error())
// 	}
// 	return &Producer{
// 		Channel: ch,
// 		logger:  r.logger,
// 		ctx:     r.ctx,
// 	}, nil
// }

/*
-------------------------------- CONSUMER METHODS --------------------------------
*/

// StartConsumers initializes and starts multiple consumers in goroutines to process messages from RabbitMQ queues.
// It adds the required number of wait groups, starts the consumers, and listens for incoming messages.
func (r *RabbitMQ) StartConsumers(
	wg *sync.WaitGroup,
	errChan chan<- error,
	stopChan <-chan struct{},
) {
	r.wg = wg
	r.errChan = errChan
	r.stopChan = stopChan

	if r.connection == nil {
		errMsg := fmt.Errorf("RabbitMQ connection is nil, ensure connection is established")
		r.errChan <- errMsg
		return
	}

	r.wg.Add(1)
	go r.startConsumer(2, USER_NOTIFICATIONS_EXCHANGE, USER_CREATE_KEY, USER_CREATION_QUEUE)
	r.wg.Add(1)
	go r.startConsumer(3, USER_NOTIFICATIONS_EXCHANGE, USER_DELETE_KEY, USER_DELETION_QUEUE)
}

// startConsumer starts a RabbitMQ consumer with specific worker.
// It sets up the consumer, declares and binds the queue, and starts listening for messages.
//
// Parameters:
//
//   - workerID: An integer representing the ID of the worker.
//
//   - exchange: The name of the RabbitMQ exchange to bind the queue to.
//
//   - routingKey: The routing key to bind the queue to the exchange.
//
//   - queueName: The name of the RabbitMQ queue to declare and bind.
//
// The function logs the progress of the consumer setup and listens for messages until
// a stop signal is received on the stopChan or an error occurs.
func (r *RabbitMQ) startConsumer(workerID int, exchange, routingKey, queueName string) {
	defer r.wg.Done()
	messageConsumer, err := r.newConsumer(workerID, exchange)
	if err != nil {
		errMsg := fmt.Errorf("[worker_%d] %s", workerID, err.Error())
		r.errChan <- errMsg
		return
	}
	r.logger.Info(fmt.Sprintf("[worker_%d] successfully declared consumer", workerID))

	// Declare queue
	if err := messageConsumer.declareQueue(queueName); err != nil {
		errMsg := fmt.Errorf("[worker_%d] %s", workerID, err.Error())
		r.errChan <- errMsg
		return
	}
	// Bind queue
	if err := messageConsumer.bindQueue(routingKey); err != nil {
		errMsg := fmt.Errorf("[worker_%d] %s", workerID, err.Error())
		r.errChan <- errMsg
		return
	}

	r.logger.Info(fmt.Sprintf("[worker_%d] start to listen...", workerID))
	messageConsumer.Listen(r.stopChan, r.errChan)

}

// Finally build and return consumer instance
// newConsumer creates a new RabbitMQ consumer with the specified worker ID and exchange name.
// It opens a new channel, declares the exchange, and returns a Consumer instance.
//
// Parameters:
//
//   - workerID: An integer representing the ID of the worker.
//
//   - exchange: A string representing the name of the exchange.
//
// Returns:
//   - *Consumer: A pointer to the newly created Consumer instance.
//   - error: An error if the channel could not be opened or the exchange could not be declared.
func (r *RabbitMQ) newConsumer(workerID int, exchange string) (*Consumer, error) {
	// create a new channel
	ch, err := r.connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel for consumer: %s", err.Error())
	}

	durable, autoDelete, internal, noWait := true, false, false, false
	if err = ch.ExchangeDeclare(
		exchange,      // name
		exchange_type, // type
		durable,       // durable
		autoDelete,    // auto-deleted
		internal,      // internal
		noWait,        // no-wait
		nil,           // arguments
	); err != nil {
		return nil, fmt.Errorf("failed to declare an exchange: %s", err.Error())
	}

	r.channels = append(r.channels, ch)
	return &Consumer{
		channel:    ch,
		ctx:        r.ctx,
		exchange:   exchange,
		logger:     r.logger,
		mongo:      r.mongo,
		workerID:   workerID,
		connection: r.connection,
	}, nil
}
