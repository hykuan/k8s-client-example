package http

import (
	"github.com/hykuan/k8s-client-example/models"
)

type apiReq interface {
	validate() error
}

type trainingReq struct {
	training models.Training
}

func (req trainingReq) validate() error {
	return req.training.Validate()
}
