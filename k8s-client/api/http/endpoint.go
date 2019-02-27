package http

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/hykuan/k8s-client-example/k8s-client"
)

func createNFSPVEndpoint(svc k8s_client.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(nfsPVReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		name, err := svc.CreateNFSPV(req.pv)
		return PVRes{name}, err
	}
}

func createPVCEndpoint(svc k8s_client.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(pvcReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		name, err := svc.CreatePVC(req.pvc)
		if err != nil {
			return nil, err
		}

		return PVCRes{name}, nil
	}
}

