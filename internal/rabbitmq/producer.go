package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Producer struct {
	// The AMQP channel used for communication with RabbitMQ.
	channel *amqp.Channel
	// The context for managing the consumer's lifecycle and cancellation.
	ctx context.Context
	// The name of the RabbitMQ exchange to bind the consumer to.
	exchange string
	//  The logger instance for logging consumer activities.
	logger *zap.Logger
	// The routing key for the producer.
	routingKey string
	// The identifier for the worker handling the consumer.
	workerID int
}

// ProduceMessage publishes a message to the queue
func (p *Producer) ProduceMessage(message string) error {
	if p.channel == nil {
		return fmt.Errorf("[worker_%d] channel is nil, ensure connection is established", p.workerID)
	}

	mandatory, immediate := false, false
	err := p.channel.PublishWithContext(
		p.ctx,        // context
		p.exchange,   // exchange
		p.routingKey, // routing key
		mandatory,    // mandatory
		immediate,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		return fmt.Errorf("[worker_%d] failed to publish message: %s", p.workerID, err.Error())
	}
	p.logger.Info(fmt.Sprintf("[worker_%d] message was produced successfully", p.workerID), zap.String("message", message))
	return nil

}
