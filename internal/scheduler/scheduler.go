package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/AnhCaooo/stormbreaker/internal/db"
	"github.com/AnhCaooo/stormbreaker/internal/electric"
	"github.com/AnhCaooo/stormbreaker/internal/helpers"
	"github.com/AnhCaooo/stormbreaker/internal/models"
	"github.com/AnhCaooo/stormbreaker/internal/rabbitmq"
	"go.uber.org/zap"
)

// Scheduler is responsible for managing and coordinating scheduled tasks.
type Scheduler struct {
	// The logger for logging RabbitMQ-related activities
	logger *zap.Logger
	// The context for managing the lifecycle of the RabbitMQ instance.
	ctx context.Context
	// The MongoDB instance for database operations.
	mongo *db.Mongo
	// Configuration settings for the RabbitMQ broker.
	brokerConfig *models.Broker
	// A pointer to a sync.WaitGroup to signal when the consumer has finished.
	wg *sync.WaitGroup
}

// NewScheduler creates a new instance of Scheduler with the provided context, logger, broker configuration, and MongoDB connection.
func NewScheduler(ctx context.Context, logger *zap.Logger, brokerConfig *models.Broker, mongo *db.Mongo) *Scheduler {
	return &Scheduler{
		logger:       logger,
		mongo:        mongo,
		ctx:          ctx,
		brokerConfig: brokerConfig,
	}
}

// StartJobs initializes and starts the scheduling jobs.
// It logs the start of the scheduling process, assigns the provided WaitGroup to the scheduler,
// increments the WaitGroup counter, and starts the PollPrice job in a new goroutine.
func (s *Scheduler) StartJobs(
	wg *sync.WaitGroup,
) {
	s.logger.Info("starting scheduling jobs...")
	s.wg = wg
	s.wg.Add(1)
	go s.PollPrice(4)
}

// StopJobs stops the scheduling jobs by decrementing the WaitGroup counter.
func (s *Scheduler) StopJobs() {
	s.logger.Info("stopping scheduling jobs...")
	s.wg.Done()
}

// PollPrice continuously polls for electricity prices
// and sends notifications if the price for the next day is available.
// The polling occurs between the specified start and end times
//
// Parameters:
//   - workerID: an integer representing the ID of the worker executing the polling job.
//
// The method performs the following steps:
//  1. Initializes a ticker to trigger every 5 seconds.
//  2. Logs the start of the polling job.
//  3. Continuously checks the current time in Helsinki.
//  4. If the current time is outside the polling hours, it pauses until the next polling period.
//  5. Checks if the price for the next day is available.
//  6. If the price is available and a notification has not been sent yet, it establishes a connection
//     to RabbitMQ, sends a notification, and then closes the connection.
//  7. Waits until the next polling period and resets the job status.
func (s *Scheduler) PollPrice(workerID int) {
	const startTime = 14
	const endTime = 22
	ticker := time.NewTicker(5 * time.Second)
	// ticker := time.NewTicker(10 * time.Minute)
	isJobDone := false
	defer ticker.Stop()

	s.logger.Info(fmt.Sprintf("[worker_%d] starting polling job...", workerID))
	for {
		currentTime, _, err := helpers.GetCurrentTimeInHelsinki()
		if err != nil {
			errMsg := fmt.Errorf("[worker_%d] failed to get current time: %s", workerID, err.Error())
			s.logger.Error(errMsg.Error())
			return
		}

		<-ticker.C
		if currentTime.Hour() < startTime || currentTime.Hour() >= endTime {
			s.logger.Info(fmt.Sprintf("[worker_%d] outside polling price hours. Pause polling until next job", workerID), zap.Time("current_time_helsinki", currentTime))
			s.waitUntilNextPollingPeriod(workerID, startTime, endTime, isJobDone)
			continue
		}

		isPriceAvailableForNotification, err := s.isTomorrowPriceAvailable(workerID)
		if err != nil {
			errMsg := fmt.Errorf("[worker_%d] failed to check if tomorrow price is available: %s", workerID, err.Error())
			s.logger.Error(errMsg.Error())
			continue
		}

		if isPriceAvailableForNotification && !isJobDone {
			s.logger.Info(fmt.Sprintf("[worker_%d] tomorrow price is available. Sending notifications...", workerID))
			rabbit := rabbitmq.NewRabbit(s.ctx, s.brokerConfig, s.logger, s.mongo)
			if err := rabbit.EstablishConnection(); err != nil {
				errMsg := fmt.Errorf("[worker_%d] failed to establish connection with RabbitMQ: %s", workerID, err.Error())
				s.logger.Error(errMsg.Error())
				return
			}
			s.logger.Info(fmt.Sprintf("[worker_%d] successfully connected to RabbitMQ", workerID))
			// send notification
			if err := rabbit.StartProducer(
				workerID,
				rabbitmq.PUSH_NOTIFICATION_EXCHANGE,
				rabbitmq.PUSH_NOTIFICATION_KEY,
				"electricity price notification",
			); err != nil {
				s.logger.Error(err.Error())
			}
			s.logger.Info(fmt.Sprintf("[worker_%d] sent notification", workerID))
			isJobDone = true
			// close connection after finish
			rabbit.CloseConnection()
			s.waitUntilNextPollingPeriod(workerID, startTime, endTime, isJobDone)
			isJobDone = false
		}
	}

}

// waitUntilNextPollingPeriod pauses the execution of a worker until the next polling period.
// It calculates the duration to wait based on the current time and the specified start and end times.
// If the job is done or the current time is past the end time, it schedules the next polling for the next day.
//
// Parameters:
//   - workerID: an integer representing the ID of the worker
//   - startTime: an integer representing the hour of the day when polling should start
//   - endTime: an integer representing the hour of the day when polling should end
//   - isJobDone: a boolean indicating whether the job is completed
//
// The function logs the duration of the wait and the time of the next polling period.
func (s Scheduler) waitUntilNextPollingPeriod(workerID, startTime, endTime int, isJobDone bool) {
	now, location, err := helpers.GetCurrentTimeInHelsinki()
	if err != nil {
		s.logger.Error(fmt.Sprintf("[worker_%d] failed to get current time", workerID), zap.Error(err))
	}
	nextStart := time.Date(
		now.Year(), now.Month(), now.Day(),
		startTime, 0, 0, 0, location,
	)

	if isJobDone || now.Hour() >= endTime {
		nextStart = nextStart.Add(24 * time.Hour) // Start polling tomorrow
	}

	duration := time.Until(nextStart)
	s.logger.Info(fmt.Sprintf("[worker_%d] on holding for %v until the next polling at %v", workerID, duration, nextStart))
	time.Sleep(duration)
}

// isTomorrowPriceAvailable checks if the price for tomorrow is available.
// It logs the process and fetches the spot price using the electric service.
// It returns a boolean indicating the availability of tomorrow's price and an error if any occurs.
//
// Parameters:
//
//	workerID (int): The ID of the worker performing the check.
//
// Returns:
//
//	bool: True if tomorrow's price is available, false otherwise.
//	error: An error object if an error occurs during the process.
func (s *Scheduler) isTomorrowPriceAvailable(workerID int) (bool, error) {
	s.logger.Info(fmt.Sprintf("[worker_%d] checking if tomorrow price is available...", workerID))
	electric := electric.NewElectric(s.logger, s.mongo, "stormbreaker")

	payloadForTodayTomorrow := electric.BuildTodayTomorrowRequestPayload()
	prices, _, err := electric.FetchSpotPrice(payloadForTodayTomorrow)
	if err != nil {
		return false, err
	}

	todayTomorrowPrice, err := helpers.MapToTodayTomorrowResponse(prices)
	if err != nil {
		return false, err
	}

	if todayTomorrowPrice.Tomorrow.Available {
		return true, nil
	}
	return false, nil
}
