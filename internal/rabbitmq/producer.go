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
	queue   *amqp.Queue
}

// DeclareQueue ensures that the queue is declared and exists before producing messages:
func (p *Producer) DeclareQueue(queueName string) error {
	if p.Channel == nil {
		return fmt.Errorf("producer channel is nil, ensure connection is established")
	}

	durable, autoDelete, exclusive, noWait := true, false, false, false
	queue, err := p.Channel.QueueDeclare(
		queueName,  // queue name
		durable,    // durable
		autoDelete, // auto-delete
		exclusive,  // exclusive
		noWait,     // no-wait
		nil,        // args
	)
	p.queue = &queue
	if err != nil {
		return fmt.Errorf("failed to declare queue: %s", err.Error())
	}
	return nil
}

// ProduceMessage publishes a message to the queue
func (p *Producer) ProduceMessage(message string) error {
	if p.Channel == nil {
		return fmt.Errorf("channel is nil, ensure connection is established")
	}

	exchange, mandatory, immediate := "", false, false
	err := p.Channel.PublishWithContext(p.ctx,
		exchange,     // exchange
		p.queue.Name, // routing key
		mandatory,    // mandatory
		immediate,    // immediate
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
