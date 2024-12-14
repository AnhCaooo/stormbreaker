package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Producer struct {
	logger  *zap.Logger
	ctx     context.Context
	Channel *amqp.Channel
}

// NewProducer receives connection client, then opens channel and build producer instance
func NewProducer(connection *amqp.Connection, logger *zap.Logger, ctx context.Context) (*Producer, error) {
	// create a new channel
	ch, err := connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %s", err.Error())
	}

	return &Producer{
		Channel: ch,
		logger:  logger,
		ctx:     ctx,
	}, nil
}

// DeclareQueue ensures that the queue is declared and exists before producing messages:
func (p *Producer) DeclareQueue(queueName string) error {
	if p.Channel == nil {
		return fmt.Errorf("producer channel is nil, ensure connection is established")
	}

	_, err := p.Channel.QueueDeclare(
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

// ProduceMessage publishes a message to the queue
func (p *Producer) ProduceMessage(queueName, message string) error {
	if p.Channel == nil {
		return fmt.Errorf("channel is nil, ensure connection is established")
	}

	err := p.Channel.PublishWithContext(p.ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		return fmt.Errorf("failed to publish message: %s", err.Error())
	}
	p.logger.Info("message was produced successfully", zap.String("message", message))
	return nil

}
