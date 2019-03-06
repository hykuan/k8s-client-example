package api

import (
	"fmt"
	"github.com/hykuan/k8s-client-example/k8s-client"
	log "github.com/hykuan/k8s-client-example/logger"
	"time"
)

var _ k8s_client.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	svc    k8s_client.Service
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc k8s_client.Service, logger log.Logger) k8s_client.Service {
	return &loggingMiddleware{logger, svc}
}

func (lm *loggingMiddleware) CreateNFSPV(nfsPV k8s_client.NFSPersistentVolume) (name string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method register for user %+v took %s to complete", nfsPV, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.CreateNFSPV(nfsPV)
}

func (lm *loggingMiddleware) CreatePVC(pvc k8s_client.PersistentVolumeClaim) (name string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method register for user %+v took %s to complete", pvc, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.CreatePVC(pvc)
}

func (lm *loggingMiddleware) CreateDeployment(deployment k8s_client.Deployment) (name string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method register for user %+v took %s to complete", deployment, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())

	return lm.svc.CreateDeployment(deployment)
}
