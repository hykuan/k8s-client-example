package grpc

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/hykuan/k8s-client-example/k8s-client"
	"golang.org/x/net/context"
)

func createNFSPVEndpoint(svc k8s_client.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createNFSPVReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		pvName, err := svc.CreateNFSPV(k8s_client.NFSPersistentVolume{
			Name: req.Name,
			Storage: req.Storage,
			Server: req.Server,
			Path: req.Path,
		})
		if err != nil {
			return createPVRes{name: "", err: err}, err
		}
		return createPVRes{name: pvName, err: nil}, nil
	}
}

func createPVCEndpoint(svc k8s_client.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createPVCReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		pvcName, err := svc.CreatePVC(k8s_client.PersistentVolumeClaim{
			Name: req.Name,
			Storage: req.Storage,
		})
		if err != nil {
			return createPVCRes{name: "", err: err}, err
		}
		return createPVCRes{name: pvcName, err: nil}, nil
	}
}

func createDeploymentEndpoint(svc k8s_client.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createDeploymentReq)
		if err := req.validate(); err != nil {
			return nil, err
		}

		resource := k8s_client.Resource{}
		if req.Resource != nil {
			resource.Memory = req.Resource.Memory
			resource.CPU = req.Resource.CPU
			resource.GPU = req.Resource.GPU
		}
		volumes := []*k8s_client.VolumeInfo{}
		for _, volume := range req.Volumes {
			volumes = append(volumes, &k8s_client.VolumeInfo{
				Name:      volume.Name,
				PVCName:   volume.PVCName,
				MountPath: volume.MountPath,
			})
		}

		deployment, err := svc.CreateDeployment(k8s_client.Deployment{
			Name:      req.Name,
			Replicas:  req.Replicas,
			Image:     req.Image,
			Resource:  &resource,
			Volumes:   volumes,
			Command:   req.Command,
			Arguments: req.Arguments,
		})
		if err != nil {
			return createDeploymentRes{name: "", err: err}, err
		}
		return createDeploymentRes{name: deployment, err: nil}, nil
	}
}
