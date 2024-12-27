package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AnhCaooo/stormbreaker/internal/db"
	"github.com/AnhCaooo/stormbreaker/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

const (
	USER_NOTIFICATIONS_EXCHANGE string = "user_notifications_exchange"
	USER_CREATE_KEY             string = "user.create"
	USER_DELETE_KEY             string = "user.delete"
	USER_CREATION_QUEUE         string = "user_creation_queue"
	USER_DELETION_QUEUE         string = "user_deletion_queue"
)

// Consumer represents a RabbitMQ consumer with necessary dependencies and configurations.
type Consumer struct {
	// The AMQP channel used for communication with RabbitMQ.
	channel *amqp.Channel
	// The AMQP connection to the RabbitMQ server.
	connection *amqp.Connection
	// The context for managing the consumer's lifecycle and cancellation.
	ctx context.Context
	// The name of the RabbitMQ exchange to bind the consumer to.
	exchange string
	//  The logger instance for logging consumer activities.
	logger *zap.Logger
	// The MongoDB instance for database operations.
	mongo *db.Mongo
	// The RabbitMQ queue to consume messages from.
	queue *amqp.Queue
	// The identifier for the worker handling the consumer.
	workerID int
}

// declareQueue declares a queue with the given name on the consumer's channel.
// It ensures the channel is not nil before attempting to declare the queue.
// If the queue declaration is successful, it assigns the declared queue to the consumer's queue field.
// Returns an error if the channel is nil or if the queue declaration fails.
func (c *Consumer) declareQueue(queueName string) error {
	if c.channel == nil {
		return fmt.Errorf("consumer channel is nil, ensure connection is established")
	}

	durable, autoDelete, exclusive, noWait := false, false, false, false
	queue, err := c.channel.QueueDeclare(
		queueName,  // queue name
		durable,    // durable
		autoDelete, // auto-delete when unused
		exclusive,  // exclusive
		noWait,     // no-wait
		nil,        // args
	)
	c.queue = &queue
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %s", err.Error())
	}
	return nil
}

// bindQueue binds the consumer's queue to the specified routing key on the exchange.
// It logs the binding action and returns an error if the binding fails.
func (c *Consumer) bindQueue(routingKey string) error {
	c.logger.Info(
		fmt.Sprintf("[worker_%d] binding queue to exchange with routing key", c.workerID),
		zap.String("queue_name", c.queue.Name),
		zap.String("exchange", c.exchange),
		zap.String("routing_key", routingKey),
	)

	if err := c.channel.QueueBind(
		c.queue.Name,
		routingKey,
		c.exchange,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to bind a queue: %s", err.Error())
	}
	return nil
}

// Listen will start to read messages from the queue
func (c *Consumer) Listen(stopChan <-chan struct{}, errChan chan<- error) {
	consumer, autoAck, exclusive, noLocal, noWait := "", false, false, false, false
	msgs, err := c.channel.Consume(
		c.queue.Name, // queue
		consumer,     // consumer
		autoAck,      // auto-ack
		exclusive,    // exclusive
		noLocal,      // no-local
		noWait,       // no-wait
		nil,          // args
	)
	if err != nil {
		errMessage := fmt.Errorf("[worker_%d] failed to register a consumer: %s", c.workerID, err.Error())
		errChan <- errMessage
		return
	}

	c.logger.Info(fmt.Sprintf("[worker_%d] waiting for messages from %s...", c.workerID, c.queue.Name))

	// Make a channel to receive messages into infinite loop.
	for {
		select {
		case <-stopChan: // Respond to shutdown signal
			c.logger.Info(fmt.Sprintf("[worker_%d] stop listening for messages from %s...", c.workerID, c.queue.Name))
			return
		case msg, ok := <-msgs:
			if !ok {
				c.logger.Info(fmt.Sprintf("[worker_%d] message channel closed", c.workerID))
				return
			}

			if err := msg.Ack(false); err != nil {
				errMsg := fmt.Errorf("[worker_%d] error acknowledging message from %s: %s", c.workerID, c.queue.Name, err.Error())
				errChan <- errMsg
				return
			}

			// Process message
			switch msg.RoutingKey {
			case USER_CREATE_KEY:
				c.logger.Info(fmt.Sprintf("[worker_%d] received a user created message", c.workerID))
				var newPriceSettingsForNewUser models.PriceSettings
				json.Unmarshal(msg.Body, &newPriceSettingsForNewUser)
				_, err := c.mongo.InsertPriceSettings(newPriceSettingsForNewUser)
				if err != nil {
					errMsg := fmt.Errorf("[worker_%d] error inserting price settings: %s", c.workerID, err.Error())
					errChan <- errMsg
				}
			case USER_DELETE_KEY:
				var deletedUserID string = string(msg.Body)
				c.logger.Info(fmt.Sprintf("[worker_%d] received a user deleted message. UserID: %s", c.workerID, deletedUserID))
				_, err := c.mongo.DeletePriceSettings(deletedUserID)
				if err != nil {
					errMsg := fmt.Errorf("[worker_%d] error delete price settings: %s", c.workerID, err.Error())
					errChan <- errMsg
				}
			default:
				c.logger.Info(fmt.Sprintf("[worker_%d] received an message from undefined routing key: %s", c.workerID, msg.RoutingKey))
			}

		}
	}
}
