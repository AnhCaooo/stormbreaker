package rabbitmq

import (
	"fmt"

	"github.com/AnhCaooo/stormbreaker/internal/constants"
	"github.com/AnhCaooo/stormbreaker/internal/db"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

const USER_SIGN_UP_NOTIFICATION = "user_sign_up_notification"

type Consumer struct {
	logger  *zap.Logger
	Channel *amqp.Channel
	queue   *amqp.Queue
	mongo   *db.Mongo
}

// DeclareQueue ensures that the queue is declared and exists before consuming messages:
func (c *Consumer) declareQueue(queueName string) error {
	if c.Channel == nil {
		return fmt.Errorf("consumer channel is nil, ensure connection is established")
	}

	durable, autoDelete, exclusive, noWait := true, false, false, false
	queue, err := c.Channel.QueueDeclare(
		queueName,  // queue name
		durable,    // durable
		autoDelete, // auto-delete
		exclusive,  // exclusive
		noWait,     // no-wait
		nil,        // args
	)
	c.queue = &queue
	if err != nil {
		return fmt.Errorf("failed to declare queue: %s", err.Error())
	}
	return nil
}

// Listen will get queue name, declare it and start to read messages from the queue
func (c *Consumer) Listen(queueName string) {
	if err := c.declareQueue(queueName); err != nil {
		c.logger.Fatal(constants.Server, zap.Error(err))
	}

	if c.Channel == nil {
		c.logger.Fatal("channel is nil, ensure connection is established")
	}

	consumer, autoAck, exclusive, noLocal, noWait := "", true, false, false, false
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
	stopChan := make(chan bool)
	go func() {
		for d := range msgs {
			c.logger.Info("Received a message", zap.Any("message", string(d.Body)))

			if err := d.Ack(false); err != nil {
				c.logger.Error("Error acknowledging message:", zap.Error(err))
			} else {
				c.logger.Info("Acknowledged message")
			}
		}
	}()
	c.logger.Info(fmt.Sprintf("[*] Waiting for messages from %s...", queueName))
	// Stop for program termination
	<-stopChan
}
