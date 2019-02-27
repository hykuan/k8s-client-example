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
	createNFSPV kitgrpc.Handler
	createPVC kitgrpc.Handler
}

// NewServer returns new K8sClientServiceServer instance.
func NewServer(svc k8s_client.Service) quai.K8SClientServiceServer {
	return &grpcServer{
		createNFSPV: kitgrpc.NewServer(
			createNFSPVEndpoint(svc),
			decodeCreateNFSPVCRequest,
			encodeCreateNFSPVCResponse,
		),
		createPVC: kitgrpc.NewServer(
			createPVCEndpoint(svc),
			decodeCreatePVCRequest,
			encodeCreatePVCResponse,
		),
	}
}

func (s *grpcServer) CreateNFSPersistentVolume(ctx context.Context, req *quai.NFSPersistentVolumeReq) (*quai.PersistentVolumeName, error) {
	_, res, err := s.createNFSPV.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}
	return res.(*quai.PersistentVolumeName), nil
}

func (s *grpcServer) CreatePersistentVolumeClaim(ctx context.Context, req *quai.PersistentVolumeClaimReq) (*quai.PersistentVolumeClaimName, error) {
	_, res, err := s.createPVC.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}
	return res.(*quai.PersistentVolumeClaimName), nil
}

func decodeCreateNFSPVCRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*quai.NFSPersistentVolumeReq)
	return createNFSPVReq{
		Name: req.Name,
		Storage: req.Storage,
		Server: req.Server,
		Path: req.Path,
	}, nil
}

func encodeCreateNFSPVCResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(createPVRes)
	return &quai.PersistentVolumeName{Value: res.name}, encodeError(res.err)
}

func decodeCreatePVCRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*quai.NFSPersistentVolumeReq)
	return createPVCReq{
		Name: req.Name,
		Storage: req.Storage,
	}, nil
}

func encodeCreatePVCResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(createPVCRes)
	return &quai.PersistentVolumeClaimName{Value: res.name}, encodeError(res.err)
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
