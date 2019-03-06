package grpc

import (
	"github.com/hykuan/k8s-client-example/k8s-client"
)

type createNFSPVReq struct {
	Name                 string
	Storage              string
	Server               string
	Path                 string
}

func (req createNFSPVReq) validate() error {
	if req.Name == "" || req.Storage == "" || req.Server == "" || req.Path == "" {
		return k8s_client.ErrMalformedEntity
	}
	return nil
}

type createPVCReq struct {
	Name                 string
	Storage              string
}

func (req createPVCReq) validate() error {
	if req.Name == "" || req.Storage == "" {
		return k8s_client.ErrMalformedEntity
	}
	return nil
}

type createDeploymentReq struct {
	deployment k8s_client.Deployment
}

func (req createDeploymentReq) validate() error {
	return req.deployment.Validate()
}
