package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Consumer struct {
	logger  *zap.Logger
	Channel *amqp.Channel
}

// NewConsumer receives connection client, then opens channel and build consumer instance
func NewConsumer(connection *amqp.Connection, logger *zap.Logger) (*Consumer, error) {
	// create a new channel
	ch, err := connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %s", err.Error())
	}

	return &Consumer{
		Channel: ch,
		logger:  logger,
	}, nil
}

// DeclareQueue ensures that the queue is declared and exists before consuming messages:
func (c *Consumer) DeclareQueue(queueName string) error {
	if c.Channel == nil {
		return fmt.Errorf("consumer channel is nil, ensure connection is established")
	}
	_, err := c.Channel.QueueDeclare(
		queueName, // queue name
		true,      // durable
		false,     // auto-delete
		false,     // exclusive
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}
	return nil
}

func (c *Consumer) ConsumeMessage(queueName string) {
	if c.Channel == nil {
		c.logger.Fatal("channel is nil, ensure connection is established")
		return
	}

	msgs, err := c.Channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		c.logger.Fatal("failed to consume message", zap.Error(err))
	}

	go func() {
		for d := range msgs {
			c.logger.Info("received message", zap.Any("message", d.Body))
		}
	}()
	c.logger.Info("[*] Waiting for messages...")
	select {} // Infinite blocking to keep the consumer running
}
