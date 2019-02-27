package http

import "github.com/hykuan/k8s-client-example/k8s-client"

type apiReq interface {
	validate() error
}

type nfsPVReq struct {
	pv k8s_client.NFSPersistentVolume
}

func (req nfsPVReq) validate() error {
	return req.pv.Validate()
}

type pvcReq struct {
	pvc k8s_client.PersistentVolumeClaim
}

func (req pvcReq) validate() error {
	return req.pvc.Validate()
}

