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
	"github.com/AnhCaooo/stormbreaker/internal/scheduler"
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
	logger := log.InitLogger(zapcore.DebugLevel)
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

	run(ctx, logger, configuration, mongo)
}

// run initializes and starts the HTTP server and RabbitMQ consumers, and listens for OS signals to gracefully shut down.
func run(ctx context.Context, logger *zap.Logger, config *models.Config, mongo *db.Mongo) {
	// Create a signal channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Error channel to listen for errors from goroutines
	var wg sync.WaitGroup
	errChan := make(chan error, 3)
	stopChan := make(chan struct{})
	// HTTP server
	httpServer := api.NewHTTPServer(ctx, logger, config, mongo)
	httpServer.Start(1, errChan, &wg)

	// RabbitMQ
	rabbitMQ := rabbitmq.NewRabbit(ctx, &config.MessageBroker, logger, mongo)
	if err := rabbitMQ.EstablishConnection(); err != nil {
		logger.Fatal("failed to establish connection with RabbitMQ", zap.Error(err))
	}
	rabbitMQ.StartConsumers(&wg, errChan, stopChan)

	// Scheduler worker
	scheduler := scheduler.NewScheduler(logger, mongo)
	scheduler.StartJobs(&wg)

	// Monitor all errors from errChan and log them
	go func() {
		for err := range errChan {
			logger.Error("error occurred", zap.Error(err))
		}
	}()

	// Wait for termination signal
	<-stop
	logger.Info("termination signal received")
	// Signal all consumers to stop
	close(stopChan)
	httpServer.Stop()
	rabbitMQ.CloseConnection()
	scheduler.StopJobs()
	// Wait for all goroutines to finish
	wg.Wait()
	// Signal all errors to stop
	close(errChan)
	logger.Info("HTTP server and RabbitMQ workers exited gracefully")
}
