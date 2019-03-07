//
// Copyright (c) 2018
// Mainflux
//
// SPDX-License-Identifier: Apache-2.0
//

package grpc

import (
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/hykuan/k8s-client-example"
	"github.com/hykuan/k8s-client-example/k8s-client"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ quai.K8SClientServiceServer = (*grpcServer)(nil)

type grpcServer struct {
	createNFSPersistentVolume   kitgrpc.Handler
	createPersistentVolumeClaim kitgrpc.Handler
	createDeployment            kitgrpc.Handler
}

// NewServer returns new K8sClientServiceServer instance.
func NewServer(svc k8s_client.Service) quai.K8SClientServiceServer {
	return &grpcServer{
		createNFSPersistentVolume: kitgrpc.NewServer(
			createNFSPVEndpoint(svc),
			decodeCreateNFSPVCRequest,
			encodeCreateNFSPVCResponse,
		),
		createPersistentVolumeClaim: kitgrpc.NewServer(
			createPVCEndpoint(svc),
			decodeCreatePVCRequest,
			encodeCreatePVCResponse,
		),
		createDeployment: kitgrpc.NewServer(
			createDeploymentEndpoint(svc),
			decodeCreateDeploymentRequest,
			encodeCreateDeploymentResponse,
		),
	}
}

func (s *grpcServer) CreateNFSPersistentVolume(ctx context.Context, req *quai.NFSPersistentVolumeReq) (*quai.PersistentVolumeName, error) {
	_, res, err := s.createNFSPersistentVolume.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}
	return res.(*quai.PersistentVolumeName), nil
}

func (s *grpcServer) CreatePersistentVolumeClaim(ctx context.Context, req *quai.PersistentVolumeClaimReq) (*quai.PersistentVolumeClaimName, error) {
	_, res, err := s.createPersistentVolumeClaim.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}
	return res.(*quai.PersistentVolumeClaimName), nil
}

func (s *grpcServer) CreateDeployment(ctx context.Context, req *quai.DeploymentReq) (*quai.DeploymentName, error) {
	_, res, err := s.createDeployment.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}
	return res.(*quai.DeploymentName), nil
}

func decodeCreateNFSPVCRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*quai.NFSPersistentVolumeReq)
	return createNFSPVReq{
		Name:    req.Name,
		Storage: req.Storage,
		Server:  req.Server,
		Path:    req.Path,
	}, nil
}

func encodeCreateNFSPVCResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(createPVRes)
	return &quai.PersistentVolumeName{Value: res.name}, encodeError(res.err)
}

func decodeCreatePVCRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*quai.PersistentVolumeClaimReq)
	return createPVCReq{
		Name:    req.Name,
		Storage: req.Storage,
	}, nil
}

func encodeCreatePVCResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(createPVCRes)
	return &quai.PersistentVolumeClaimName{Value: res.name}, encodeError(res.err)
}

func decodeCreateDeploymentRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*quai.DeploymentReq)

	resource := Resource{}
	if req.Resource != nil {
		resource.Memory = req.Resource.Memory
		resource.CPU = req.Resource.CPU
		resource.GPU = req.Resource.GPU
	}
	volumes := []*VolumeInfo{}
	for _, volume := range req.Volumes {
		volumes = append(volumes, &VolumeInfo{
			Name:      volume.Name,
			PVCName:   volume.PVCName,
			MountPath: volume.MountPath,
		})
	}

	return createDeploymentReq{
			Name:      req.Name,
			Replicas:  req.Replicas,
			Image:     req.Image,
			Resource:  &resource,
			Volumes:   volumes,
			Command:   req.Command,
			Arguments: req.Arguments,
	}, nil
}

func encodeCreateDeploymentResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(createDeploymentRes)
	return &quai.DeploymentName{Value: res.name}, encodeError(res.err)
}

func encodeError(err error) error {
	if err == nil {
		return nil
	}

	switch err {
	case k8s_client.ErrMalformedEntity:
		return status.Error(codes.InvalidArgument, "received invalid token request")
	case k8s_client.ErrUnauthorizedAccess:
		return status.Error(codes.Unauthenticated, "failed to identify user from token")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
