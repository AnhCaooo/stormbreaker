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

func (r *Rabbit) EstablishConnection() (*amqp.Connection, error) {
	connection, err := amqp.Dial(r.getURI())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %s", err.Error())
	}

	r.connection = connection
	r.logger.Info("Successfully connected to RabbitMQ")
	return connection, nil
}

func (r *Rabbit) getURI() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", r.config.Username, r.config.Password, r.config.Host, r.config.Port)
}
