package api

import (
	"github.com/hykuan/k8s-client-example/k8s-client"
	"time"

	"github.com/go-kit/kit/metrics"
)

var _ k8s_client.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     k8s_client.Service
}

// MetricsMiddleware instruments core service by tracking request count and
// latency.
func MetricsMiddleware(svc k8s_client.Service, counter metrics.Counter, latency metrics.Histogram) k8s_client.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (ms *metricsMiddleware) CreateNFSPV(nfsPV k8s_client.NFSPersistentVolume) (name string, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "register").Add(1)
		ms.latency.With("method", "register").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.CreateNFSPV(nfsPV)
}

func (ms *metricsMiddleware) CreatePVC(pvc k8s_client.PersistentVolumeClaim) (name string, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "login").Add(1)
		ms.latency.With("method", "login").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.CreatePVC(pvc)
}