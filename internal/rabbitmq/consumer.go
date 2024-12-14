package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Consumer struct {
	logger  *zap.Logger
	channel *amqp.Channel
}

func NewConsumer(channel *amqp.Channel, logger *zap.Logger) *Consumer {
	return &Consumer{
		channel: channel,
		logger:  logger,
	}
}

// DeclareQueue ensures that the queue is declared and exists before consuming messages:
func (c *Consumer) DeclareQueue(queueName string) error {
	_, err := c.channel.QueueDeclare(
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
	if c.channel == nil {
		c.logger.Fatal("channel is nil, ensure connection is established")
		return
	}

	msgs, err := c.channel.Consume(
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
