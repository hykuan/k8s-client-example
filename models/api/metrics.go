package api

import (
	"time"

	"github.com/go-kit/kit/metrics"

	"github.com/hykuan/k8s-client-example/models"
)

var _ models.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     models.Service
}

// MetricsMiddleware instruments core service by tracking request count and
// latency.
func MetricsMiddleware(svc models.Service, counter metrics.Counter, latency metrics.Histogram) models.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (ms *metricsMiddleware) StartTraining(training models.Training) (name string, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "register").Add(1)
		ms.latency.With("method", "register").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.StartTraining(training)
}
