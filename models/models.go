package models

type MountedPersistentVolumeClaim struct {
	PVCName   string
	MountPath string
}

type Training struct {
	Name      string
	Image     string
	DataSet   *MountedPersistentVolumeClaim
	Model     *MountedPersistentVolumeClaim
	GPU       uint64
	Command   []string
	Arguments []string
}

func (t Training) Validate() error {
	if t.Name == "" || t.Image == "" || t.DataSet == nil || t.Model == nil {
		return ErrMalformedEntity
	}

	if t.DataSet.MountPath == "" || t.DataSet.PVCName == "" || t.Model.PVCName == "" || t.Model.MountPath == "" {
		return ErrMalformedEntity
	}

	return nil
}
