package grpc

import "github.com/hykuan/k8s-client-example/models"

type trainingReq struct {
	training models.Training
}

func (req trainingReq) validate() error {
	return req.training.Validate()
}
