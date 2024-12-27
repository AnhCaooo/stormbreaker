package rabbitmq

import (
	"context"
	"fmt"

	"github.com/AnhCaooo/stormbreaker/internal/db"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

const (
	USER_NOTIFICATIONS_EXCHANGE string = "user_notifications_exchange"
	USER_TEST_KEY1              string = "user.test1" // todo: to be removed
	USER_TEST_KEY2              string = "user.test2" // todo: to be removed
	USER_CREATED_KEY            string = "user.created"
	USER_DELETED_KEY            string = "user.deleted"
	USER_TEST_QUEUE1            string = "user_test_queue1" // todo: to be removed
	USER_TEST_QUEUE2            string = "user_test_queue2" // todo: to be removed
	USER_CREATED_QUEUE          string = "user_created_queue"
	USER_DELETED_QUEUE          string = "user_deleted_queue"
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
			c.logger.Info(
				fmt.Sprintf("[worker_%d] received a message from %s", c.workerID, c.queue.Name),
				zap.Any("message", string(msg.Body)),
			)
			// Process message
			if err := msg.Ack(false); err != nil {
				errMsg := fmt.Errorf("[worker_%d] error acknowledging message from %s: %s", c.workerID, c.queue.Name, err.Error())
				errChan <- errMsg
				return
			} else {
				c.logger.Info(
					fmt.Sprintf("[worker_%d] acknowledged message from %s", c.workerID, c.queue.Name),
				)
			}
		}
	}
}
