package rabbitmq

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/AnhCaooo/stormbreaker/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Rabbit struct {
	config     *models.Broker
	ctx        context.Context
	logger     *zap.Logger
	connection *amqp.Connection
}

func NewRabbit(ctx context.Context, config *models.Broker, logger *zap.Logger) *Rabbit {
	return &Rabbit{
		ctx:    ctx,
		config: config,
		logger: logger,
	}
}

// EstablishConnection tries to establish a connection with RabbitMQ server
func (r *Rabbit) EstablishConnection() (*amqp.Connection, error) {
	connection, err := amqp.Dial(r.getURI())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %s", err.Error())
	}

	r.connection = connection
	r.logger.Info("Successfully connected to RabbitMQ")
	return connection, nil
}

// NewProducer retrieves connection client, then opens channel and build producer instance
func (r *Rabbit) NewProducer() (*Producer, error) {
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

// NewConsumer retrieves connection client, then opens channel and build consumer instance
func (r *Rabbit) NewConsumer() (*Consumer, error) {
	// create a new channel
	ch, err := r.connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel for consumer: %s", err.Error())
	}

	return &Consumer{
		Channel: ch,
		logger:  r.logger,
	}, nil
}

func (r *Rabbit) getURI() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", r.config.Username, r.config.Password, r.config.Host, r.config.Port)
}
