package http

import (
	"net/http"

	"github.com/hykuan/k8s-client-example"
)

var (
	_ quai.Response = (*TrainingRes)(nil)
)

type TrainingRes struct {
	Name string `json:"name,omitempty"`
}

func (res TrainingRes) Code() int {
	return http.StatusCreated
}

func (res TrainingRes) Headers() map[string]string {
	return map[string]string{}
}

func (res TrainingRes) Empty() bool {
	return res.Name == ""
}
