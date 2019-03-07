package grpc

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"

	"github.com/hykuan/k8s-client-example/models"
)

func startTrainingEndpoint(svc models.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(trainingReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		trainingName, err := svc.StartTraining(models.Training{
			Name:  req.training.Name,
			Image: req.training.Image,
			DataSet: &models.MountedPersistentVolumeClaim{
				PVCName:   req.training.DataSet.PVCName,
				MountPath: req.training.DataSet.MountPath,
			},
			Model: &models.MountedPersistentVolumeClaim{
				PVCName:   req.training.Model.PVCName,
				MountPath: req.training.Model.MountPath,
			},
			GPU:       req.training.GPU,
			Command:   req.training.Command,
			Arguments: req.training.Arguments,
		})
		if err != nil {
			return trainingRes{name: "", err: err}, err
		}
		return trainingRes{name: trainingName, err: nil}, nil
	}
}
