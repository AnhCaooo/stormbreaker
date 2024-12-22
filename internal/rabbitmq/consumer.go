package rabbitmq

import (
	"fmt"

	"github.com/AnhCaooo/stormbreaker/internal/db"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

const (
	USER_NOTIFICATIONS_EXCHANGE string = "user_notifications_exchange"
	USER_CREATED_KEY            string = "user.created"
	USER_DELETED_KEY            string = "user.deleted"
	USER_CREATED_QUEUE          string = "user_created_queue"
	USER_DELETED_QUEUE          string = "user_deleted_queue"
)

type Consumer struct {
	Channel  *amqp.Channel
	exchange string
	logger   *zap.Logger
	mongo    *db.Mongo
	queue    *amqp.Queue
}

// DeclareQueue ensures that the queue is declared and exists before consuming messages:
func (c *Consumer) declareQueue(queueName string) error {
	if c.Channel == nil {
		return fmt.Errorf("consumer channel is nil, ensure connection is established")
	}

	durable, autoDelete, exclusive, noWait := false, false, false, false
	queue, err := c.Channel.QueueDeclare(
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

func (c *Consumer) bindQueue(routingKey string) error {
	c.logger.Info(fmt.Sprintf("Binding '%s' to '%s' with routing key '%s'",
		c.queue.Name, c.exchange, routingKey))
	if err := c.Channel.QueueBind(
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
func (c *Consumer) Listen() {
	consumer, autoAck, exclusive, noLocal, noWait := "", false, false, false, false
	msgs, err := c.Channel.Consume(
		c.queue.Name, // queue
		consumer,     // consumer
		autoAck,      // auto-ack
		exclusive,    // exclusive
		noLocal,      // no-local
		noWait,       // no-wait
		nil,          // args
	)
	if err != nil {
		c.logger.Fatal("failed to register a consumer:", zap.Error(err))
	}

	c.logger.Info(fmt.Sprintf("[*] Waiting for messages from %s...", c.queue.Name))

	// Make a channel to receive messages into infinite loop.
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			c.logger.Info(
				fmt.Sprintf("[*] Received a message from %s", c.queue.Name),
				zap.Any("message", string(d.Body)),
			)

			if err := d.Ack(false); err != nil {
				c.logger.Error(
					fmt.Sprintf("[*] Error acknowledging message from %s:", c.queue.Name),
					zap.Error(err),
				)
			} else {
				c.logger.Info(
					fmt.Sprintf("[*] Acknowledged message from %s", c.queue.Name),
				)
			}
		}
	}()
	// Stop for program termination
	<-forever
}
