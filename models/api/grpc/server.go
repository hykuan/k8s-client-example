//
// Copyright (c) 2018
// Mainflux
//
// SPDX-License-Identifier: Apache-2.0
//

package grpc

import (
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hykuan/k8s-client-example"
	"github.com/hykuan/k8s-client-example/models"
)

var _ quai.ModelServiceServer = (*grpcServer)(nil)

type grpcServer struct {
	startTraining kitgrpc.Handler
}

// NewServer returns new ModelServiceServer instance.
func NewServer(svc models.Service) quai.ModelServiceServer {
	return &grpcServer{
		startTraining: kitgrpc.NewServer(
			startTrainingEndpoint(svc),
			decodeTrainingRequest,
			encodeTrainingResponse,
		),
	}
}

func (s *grpcServer) StartTraining(ctx context.Context, req *quai.TrainingReq) (*quai.Training, error) {
	_, res, err := s.startTraining.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}
	return res.(*quai.Training), nil
}

func decodeTrainingRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*quai.TrainingReq)

	dataSet := models.MountedPersistentVolumeClaim{}
	if req.DataSet != nil {
		dataSet.PVCName = req.DataSet.PVCName
		dataSet.MountPath = req.DataSet.MountPath
	}

	model := models.MountedPersistentVolumeClaim{}
	if req.Model != nil {
		model.PVCName = req.Model.PVCName
		model.MountPath = req.Model.MountPath
	}

	return trainingReq{
		training: models.Training{
			Name:      req.Name,
			Image:     req.Image,
			DataSet:   &dataSet,
			Model:     &model,
			GPU:       req.GPU,
			Command:   req.Command,
			Arguments: req.Arguments,
		},
	}, nil
}

func encodeTrainingResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(trainingRes)
	return &quai.Training{Value: res.name}, encodeError(res.err)
}

func encodeError(err error) error {
	if err == nil {
		return nil
	}

	switch err {
	case models.ErrMalformedEntity:
		return status.Error(codes.InvalidArgument, "received invalid token request")
	case models.ErrUnauthorizedAccess:
		return status.Error(codes.Unauthenticated, "failed to identify user from token")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
