package rabbitmq

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/AnhCaooo/stormbreaker/internal/db"
	"github.com/AnhCaooo/stormbreaker/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	reconnectDelay    = 5 * time.Second
	reconnectAttempts = 5
)

type RabbitMQ struct {
	config     *models.Broker
	connection *amqp.Connection
	ctx        context.Context
	logger     *zap.Logger
	mongo      *db.Mongo
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
func (r *RabbitMQ) establishConnection() (err error) {
	r.connection, err = amqp.Dial(r.getURI())
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %s", err.Error())
	}
	return nil
}

func (r *RabbitMQ) getURI() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", r.config.Username, r.config.Password, "localhost", r.config.Port)
}

// monitorConnection creates a go channel and a goroutine to monitor the connection.
// if connection is lost, then reconnect
func (r *RabbitMQ) monitorConnection() {
	notifyClose := make(chan *amqp.Error)
	r.connection.NotifyClose(notifyClose)

	for {
		err := <-notifyClose
		if err != nil {
			r.logger.Warn("Connection closed. Reconnecting...", zap.Error(err))
			var newConn *amqp.Connection
			var reconnectErr error
			for {
				if reconnectErr = r.establishConnection(); reconnectErr == nil {
					r.logger.Info("Reconnected to RabbitMQ in goroutine")
					r.connection = newConn
					notifyClose = make(chan *amqp.Error)
					r.connection.NotifyClose(notifyClose)
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
	if err := r.establishConnection(); err != nil {
		return nil, err
	}

	go r.monitorConnection()
	// create a new channel
	ch, err := r.connection.Channel()
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

// StartConsumer create new RabbitMQ consumer based on given queue name.
// Then listen to incoming messages from the queue
func (r *RabbitMQ) StartConsumer(
	workerID int,
	wg *sync.WaitGroup,
	errChan chan<- error,
	stopChan <-chan struct{},
	exchange, routingKey, queueName string,
) {
	defer wg.Done()
	messageConsumer, err := r.newConsumer(workerID, exchange)
	if err != nil {
		errMsg := fmt.Errorf("[* worker %d] %s", workerID, err.Error())
		errChan <- errMsg
		return
	}
	r.logger.Info(fmt.Sprintf("[* worker %d] Successfully connected to RabbitMQ and declared consumer", workerID))

	// Declare queue
	if err := messageConsumer.declareQueue(queueName); err != nil {
		errMsg := fmt.Errorf("[* worker %d] %s", workerID, err.Error())
		errChan <- errMsg
		return
	}
	// Bind queue
	if err := messageConsumer.bindQueue(routingKey); err != nil {
		errMsg := fmt.Errorf("[* worker %d] %s", workerID, err.Error())
		errChan <- errMsg
		return
	}

	r.logger.Info(fmt.Sprintf("[* worker %d] start to listen...", workerID))
	messageConsumer.Listen(stopChan, errChan)

}

// NewConsumer retrieves connection client, then opens new channel,
// and declare exchange.
// Finally build and return consumer instance
func (r *RabbitMQ) newConsumer(workerID int, exchange string) (*Consumer, error) {
	exchange_type := "topic"
	if err := r.establishConnection(); err != nil {
		return nil, err
	}
	// go r.monitorConnection()
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
