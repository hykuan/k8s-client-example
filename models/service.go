//
// Copyright (c) 2018
// Mainflux
//
// SPDX-License-Identifier: Apache-2.0
//

package models

import (
	"errors"
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
)

// Service specifies an API that must be fullfiled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	StartTraining(req quai.TrainingReq) error
}

var _ Service = (*modelsService)(nil)

type modelsService struct {
	//users  UserRepository
	//hasher Hasher
	//idp    IdentityProvider
}

// New instantiates the users service implementation.
func New() Service {
	return &modelsService{}
}

func (svc modelsService) StartTraining(req quai.TrainingReq) error {
	return nil
}
