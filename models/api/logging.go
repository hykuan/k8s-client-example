package api

import (
	"fmt"
	"time"

	log "github.com/hykuan/k8s-client-example/logger"
	"github.com/hykuan/k8s-client-example/models"
)

var _ models.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	svc    models.Service
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc models.Service, logger log.Logger) models.Service {
	return &loggingMiddleware{logger, svc}
}

func (lm *loggingMiddleware) StartTraining(training models.Training) (name string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method register for user %+v took %s to complete", training, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.StartTraining(training)
}
