package http

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/hykuan/k8s-client-example/models"
)

func startTrainingEndpoint(svc models.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(trainingReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		name, err := svc.StartTraining(req.training)
		return TrainingRes{name}, err
	}
}
