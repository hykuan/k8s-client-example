package http

import (
	"github.com/hykuan/k8s-client-example"
	"net/http"
)

var (
	_ quai.Response = (*PVRes)(nil)
	_ quai.Response = (*PVCRes)(nil)
)

type PVRes struct {
	Name string `json:"name,omitempty"`
}

func (res PVRes) Code() int {
	return http.StatusCreated
}

func (res PVRes) Headers() map[string]string {
	return map[string]string{}
}

func (res PVRes) Empty() bool {
	return res.Name == ""
}

type PVCRes struct {
	Name string `json:"name,omitempty"`
}

func (res PVCRes) Code() int {
	return http.StatusCreated
}

func (res PVCRes) Headers() map[string]string {
	return map[string]string{}
}

func (res PVCRes) Empty() bool {
	return res.Name == ""
}
