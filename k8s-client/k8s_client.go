package k8s_client

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	defaultReplicas = 1
)

type NFSPersistentVolume struct {
	Name    string
	Storage string
	Server  string
	Path    string
}

func (pv NFSPersistentVolume) Validate() error {
	if pv.Name == "" || pv.Storage == "" || pv.Server == "" || pv.Path == "" {
		return ErrMalformedEntity
	}

	return nil
}

type PersistentVolumeClaim struct {
	Name    string
	Storage string
}

func (pvc PersistentVolumeClaim) Validate() error {
	if pvc.Name == "" || pvc.Storage == "" {
		return ErrMalformedEntity
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

type Deployment struct {
	Name      string
	Replicas  int32
	Image     string
	Resource  *Resource
	Volumes   []*VolumeInfo
	Command   []string
	Arguments []string
}

func (d Deployment) Validate() error {
	if d.Name == "" || d.Image == "" {
		return ErrMalformedEntity
	}

	return nil
}

func (d *Deployment) AssignDefaultValue() {
	if d.Replicas <= 0 {
		d.Replicas = defaultReplicas
	}
}

func (d Deployment) GetResourceList() v1.ResourceList {
	list := v1.ResourceList{}

	if d.Resource == nil {
		return list
	}

	if cpu, err := resource.ParseQuantity(d.Resource.CPU); err == nil {
		list["cpu"] = cpu
	}

	if memory, err := resource.ParseQuantity(d.Resource.Memory); err == nil {
		list["memory"] = memory
	}

	if gpu, err := resource.ParseQuantity(d.Resource.GPU); err == nil {
		list["nvidia.com/gpu"] = gpu
	}

	return list
}

func (d Deployment) GetVolumes() []v1.Volume {
	var volumes []v1.Volume
	for _, v := range d.Volumes {
		volumes = append(volumes, v1.Volume{
			Name: v.Name,
			VolumeSource: v1.VolumeSource{
				PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
					ClaimName: v.Name,
				},
			},
		})
	}

	return volumes
}

func (d Deployment) GetVolumeMounts() []v1.VolumeMount {
	var volumeMounts []v1.VolumeMount
	for _, v := range d.Volumes {
		volumeMounts = append(volumeMounts, v1.VolumeMount{
			Name: v.Name,
			MountPath: v.MountPath,
		})
	}

	return volumeMounts
}

