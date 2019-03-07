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

type Resource struct {
	CPU    string
	Memory string
	GPU    string
}

type VolumeInfo struct {
	Name      string
	PVCName   string
	MountPath string
}

type createDeploymentReq struct {
	Name      string
	Replicas  int32
	Image     string
	Resource  *Resource
	Volumes   []*VolumeInfo
	Command   []string
	Arguments []string
}

func (req createDeploymentReq) validate() error {
	if req.Name == "" || req.Image == "" {
		return k8s_client.ErrMalformedEntity
	}

	return nil
}
