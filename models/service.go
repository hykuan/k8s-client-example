//
// Copyright (c) 2018
// Mainflux
//
// SPDX-License-Identifier: Apache-2.0
//

package models

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/hykuan/k8s-client-example"
)

var (
	// ErrConflict indicates usage of the existing email during account
	// registration.
	ErrConflict = errors.New("email already taken")

	// ErrMalformedEntity indicates malformed entity specification (e.g.
	// invalid username or password).
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent entity")

	// ErrNotFound indicates a non-existent entity request.
	ErrK8SCreateDeployment = errors.New("create deployment failed")
)

// Service specifies an API that must be fullfiled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	StartTraining(req Training) (string, error)
}

var _ Service = (*modelsService)(nil)

type modelsService struct {
	k8s quai.K8SClientServiceClient
}

// New instantiates the users service implementation.
func New(k8sClient quai.K8SClientServiceClient) Service {
	return &modelsService{
		k8s: k8sClient,
	}
}

func (svc *modelsService) StartTraining(training Training) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	deployment, err := svc.k8s.CreateDeployment(ctx, &quai.DeploymentReq{
		Name:      training.Name,
		Image:     training.Image,
		Command:   training.Command,
		Arguments: training.Arguments,
		Resource: &quai.Resource{
			GPU: strconv.FormatUint(training.GPU, 10),
		},
		Volumes: []*quai.VolumeInfo{
			{Name: training.DataSet.PVCName, PVCName: training.DataSet.PVCName, MountPath: training.DataSet.MountPath},
			{Name: training.Model.PVCName, PVCName: training.Model.PVCName, MountPath: training.Model.MountPath},
		},
	})

	if err != nil {
		return "", ErrK8SCreateDeployment
	}

	return deployment.Value, err
}
