package scheduler

import (
	"fmt"
	"sync"
	"time"

	"github.com/AnhCaooo/stormbreaker/internal/db"
	"github.com/AnhCaooo/stormbreaker/internal/electric"
	"github.com/AnhCaooo/stormbreaker/internal/helpers"
	"go.uber.org/zap"
)

type Scheduler struct {
	logger *zap.Logger
	mongo  *db.Mongo
	// A pointer to a sync.WaitGroup to signal when the consumer has finished.
	wg *sync.WaitGroup
}

func NewScheduler(logger *zap.Logger, mongo *db.Mongo) *Scheduler {
	return &Scheduler{
		logger: logger,
		mongo:  mongo,
	}
}

func (s *Scheduler) StartJobs(
	wg *sync.WaitGroup,
) {
	s.logger.Info("starting scheduling jobs...")
	s.wg = wg
	s.wg.Add(1)
	go s.PollPrice(4)
}

func (s *Scheduler) StopJobs() {
	s.logger.Info("stopping scheduling jobs...")
	s.wg.Done()
}

func (s *Scheduler) PollPrice(workerID int) {
	const startTime = 14
	const endTime = 17
	ticker := time.NewTicker(5 * time.Second)
	// ticker := time.NewTicker(10 * time.Minute)

	defer ticker.Stop()

	for {
		currentTime, err := helpers.GetCurrentTimeInHelsinki()
		if err != nil {
			errMsg := fmt.Errorf("[worker_%d] failed to get current time: %s", workerID, err.Error())
			s.logger.Error(errMsg.Error())
			return
		}

		<-ticker.C
		s.logger.Info(fmt.Sprintf("[worker_%d] starting polling job...", workerID))

		if currentTime.Hour() < startTime || currentTime.Hour() >= endTime {
			s.logger.Info(fmt.Sprintf("[worker_%d] outside polling hours. Pause polling until next job", workerID), zap.Time("current_time_helsinki", currentTime))
			s.waitUntilNextPollingPeriod(workerID, startTime, endTime)
		}

		isPriceAvailableForNotification, err := s.isTomorrowPriceAvailable(workerID)
		if err != nil {
			errMsg := fmt.Errorf("[worker_%d] failed to check if tomorrow price is available: %s", workerID, err.Error())
			s.logger.Error(errMsg.Error())
			return
		}
		if isPriceAvailableForNotification {
			s.logger.Info(fmt.Sprintf("[worker_%d] tomorrow price is available. Sending notification", workerID))
			// create new rabbitmq connection
			// send notification
			// close connection after finish
			return
		}
	}

}

func (s Scheduler) waitUntilNextPollingPeriod(workerID, startTime, endTime int) {
	nowInUTC, err := helpers.GetCurrentTimeInHelsinki()
	if err != nil {
		s.logger.Error(fmt.Sprintf("[worker_%d] failed to get current time", workerID), zap.Error(err))
	}

	nextStart := time.Date(
		nowInUTC.Year(), nowInUTC.Month(), nowInUTC.Day(),
		startTime, 0, 0, 0, nowInUTC.Location(),
	)
	if nowInUTC.Hour() >= endTime {
		nextStart = nextStart.Add(24 * time.Hour) // Start polling tomorrow
		s.logger.Info(fmt.Sprintf("[worker_%d] It's past %d:00. Starting polling at %d:00 tomorrow.", workerID, endTime, startTime))
	}
	duration := time.Until(nextStart)
	s.logger.Info(fmt.Sprintf("[worker_%d] Waiting for %v until the next polling period", workerID, duration))
	time.Sleep(duration)
}

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
