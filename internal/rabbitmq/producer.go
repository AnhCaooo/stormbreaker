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
	channel *amqp.Channel
}

func NewProducer(channel *amqp.Channel, logger *zap.Logger, ctx context.Context) *Producer {
	return &Producer{
		channel: channel,
		logger:  logger,
		ctx:     ctx,
	}
}

// DeclareQueue ensures that the queue is declared and exists before producing messages:
func (p *Producer) DeclareQueue(queueName string) error {
	_, err := p.channel.QueueDeclare(
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
	if p.channel == nil {
		return fmt.Errorf("channel is nil, ensure connection is established")
	}

	err := p.channel.PublishWithContext(p.ctx,
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
