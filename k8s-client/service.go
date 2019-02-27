//
// Copyright (c) 2018
// Mainflux
//
// SPDX-License-Identifier: Apache-2.0
//

package k8s_client

import (
	"errors"
	"github.com/kubernetes/client-go/kubernetes/typed/core/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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
	CreateNFSPV(nfsPV NFSPersistentVolume) (string, error)
	CreatePVC(pvc PersistentVolumeClaim) (string, error)
}

var _ Service = (*k8sClientService)(nil)

type k8sClientService struct {
	clientSet *kubernetes.Clientset
	pvClient  v1.PersistentVolumeInterface
	pvcClient  v1.PersistentVolumeClaimInterface
}

// New instantiates the users service implementation.
func New(clientSet *kubernetes.Clientset) Service {
	return &k8sClientService{
		clientSet: clientSet,
		pvClient:  clientSet.CoreV1().PersistentVolumes(),
		pvcClient:  clientSet.CoreV1().PersistentVolumeClaims(apiv1.NamespaceDefault),
	}
}

func (svc k8sClientService) CreateNFSPV(nfsPV NFSPersistentVolume) (string, error) {
	storage, err := resource.ParseQuantity(nfsPV.Storage)
	if err != nil {
		return "", err
	}

	pv, err := svc.pvClient.Create(&apiv1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: nfsPV.Name,
		},
		Spec: apiv1.PersistentVolumeSpec{
			Capacity: apiv1.ResourceList{
				"storage": storage,
			},
			AccessModes: []apiv1.PersistentVolumeAccessMode{
				apiv1.ReadWriteOnce,
			},
			PersistentVolumeSource: apiv1.PersistentVolumeSource{
				NFS: &apiv1.NFSVolumeSource{
					Server: nfsPV.Server,
					Path: nfsPV.Path,
				},
			},
		},
	})

	if err != nil {
		return "", err
	}

	return pv.Name, nil
}

func (svc k8sClientService) CreatePVC(pvc PersistentVolumeClaim) (string, error) {
	volumeMode := apiv1.PersistentVolumeBlock
	storage, err := resource.ParseQuantity(pvc.Storage)
	if err != nil {
		return "", err
	}

	pvClaim, err := svc.pvcClient.Create(&apiv1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: pvc.Name,
		},
		Spec: apiv1.PersistentVolumeClaimSpec{
			AccessModes: []apiv1.PersistentVolumeAccessMode{
				apiv1.ReadWriteOnce,
			},
			VolumeMode: &volumeMode,
			Resources: apiv1.ResourceRequirements{
				Requests: apiv1.ResourceList{
					"storage": storage,
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	return pvClaim.Name, nil
}
