package k8s_client

type NFSPersistentVolume struct {
	Name string
	Storage string
	Server string
	Path string
}

func (pv NFSPersistentVolume) Validate() error {
	if pv.Name == "" || pv.Storage == "" || pv.Server == "" || pv.Path == "" {
		return ErrMalformedEntity
	}

	return nil
}

type PersistentVolumeClaim struct {
	Name string
	Storage string
}

func (pvc PersistentVolumeClaim) Validate() error {
	if pvc.Name == "" || pvc.Storage == "" {
		return ErrMalformedEntity
	}

	return nil
}

type Deployment struct {
	Name string
	Image string
}

func (d Deployment) Validate() error {
	if d.Name == "" || d.Image == "" {
		return ErrMalformedEntity
	}

	return nil
}
