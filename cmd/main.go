// AnhCao 2024
package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/AnhCaooo/go-goods/log"
	_ "github.com/AnhCaooo/stormbreaker/docs"
	"github.com/AnhCaooo/stormbreaker/internal/api"
	"github.com/AnhCaooo/stormbreaker/internal/config"
	"github.com/AnhCaooo/stormbreaker/internal/constants"
	"github.com/AnhCaooo/stormbreaker/internal/db"
	"github.com/AnhCaooo/stormbreaker/internal/models"
	"github.com/AnhCaooo/stormbreaker/internal/rabbitmq"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//	@title			Stormbreaker API (electric service)
//	@version		1.0.0
//	@description	Service for retrieving information about market electric price in Finland.

//	@contact.name	Anh Cao
//	@contact.email	anhcao4922@gmail.com

// @host		localhost:5001
// @BasePath	/
func main() {
	ctx := context.Background()

	// todo: implement to accept a dynamic log level
	// Initialize logger
	logger := log.InitLogger(zapcore.InfoLevel)
	defer logger.Sync()

	configuration := &models.Config{}
	// Load and validate configuration
	err := config.LoadFile(configuration)
	if err != nil {
		logger.Fatal(constants.Server, zap.Error(err))
	}

	// Initialize database connection
	mongo := db.NewMongo(ctx, &configuration.Database, logger)
	if err := mongo.EstablishConnection(); err != nil {
		logger.Fatal(constants.Server, zap.Error(err))
	}
	defer mongo.Client.Disconnect(ctx)

	// Start server
	run(ctx, logger, configuration, mongo)
}

func run(ctx context.Context, logger *zap.Logger, config *models.Config, mongo *db.Mongo) {
	// Create a signal channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Error channel to listen for errors from goroutines
	var wg sync.WaitGroup
	errChan := make(chan error, 3)
	stopChan := make(chan struct{})

	// Run the server in a separate goroutine
	httpServer := api.NewHTTPServer(ctx, logger, config, mongo)
	wg.Add(1)
	go httpServer.Start(1, errChan, &wg)

	rabbitMQ := rabbitmq.NewRabbit(ctx, &config.MessageBroker, logger, mongo)
	if err := rabbitMQ.EstablishConnection(); err != nil {
		logger.Fatal("Failed to establish connection with RabbitMQ", zap.Error(err))
	}
	// Initialize RabbitMQ connections in a separate goroutine
	wg.Add(1)
	go func() {
		rabbitMQ.StartConsumer(
			2,
			&wg,
			errChan,
			stopChan,
			rabbitmq.USER_NOTIFICATIONS_EXCHANGE,
			rabbitmq.USER_TEST_KEY1,
			rabbitmq.USER_TEST_QUEUE1,
		)
	}()
	wg.Add(1)
	go func() {
		rabbitMQ.StartConsumer(
			3,
			&wg,
			errChan,
			stopChan,
			rabbitmq.USER_NOTIFICATIONS_EXCHANGE,
			rabbitmq.USER_TEST_KEY2,
			rabbitmq.USER_TEST_QUEUE2,
		)
	}()

	// Monitor all errors from errChan and log them
	go func() {
		for err := range errChan {
			logger.Error("Error occurred", zap.Error(err))
		}
	}()

	// Wait for termination signal
	<-stop
	logger.Info("Termination signal received")
	// Signal all consumers to stop
	close(stopChan)
	httpServer.Stop()
	rabbitMQ.CloseConnection()
	// Wait for all goroutines to finish
	wg.Wait()
	// Signal all errors to stop
	close(errChan)
	logger.Info("Server and RabbitMQ workers exited gracefully")
}
