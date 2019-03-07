package grpc

import (
	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/hykuan/k8s-client-example"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var _ quai.K8SClientServiceClient = (*grpcClient)(nil)

type grpcClient struct {
	createNFSPersistentVolume   endpoint.Endpoint
	createPersistentVolumeClaim endpoint.Endpoint
	createDeployment            endpoint.Endpoint
}

// NewClient returns new gRPC client instance.
func NewClient(conn *grpc.ClientConn) quai.K8SClientServiceClient {
	svcName := "quai.K8sClientService"

	return &grpcClient{
		createNFSPersistentVolume: kitgrpc.NewClient(
			conn,
			svcName,
			"CreateNFSPersistentVolume",
			encodeCreateNFSPVRequest,
			decodeCreateNFSPVResponse,
			quai.PersistentVolumeName{},
		).Endpoint(),
		createPersistentVolumeClaim: kitgrpc.NewClient(
			conn,
			svcName,
			"CreatePersistentVolumeClaim",
			encodeCreatePVCRequest,
			decodeCreatePVCResponse,
			quai.PersistentVolumeClaimName{},
		).Endpoint(),
		createDeployment: kitgrpc.NewClient(
			conn,
			svcName,
			"CreateDeployment",
			encodeCreateDeploymentRequest,
			decodeCreateDeploymentResponse,
			quai.DeploymentName{},
		).Endpoint(),
	}
}

func (client *grpcClient) CreateNFSPersistentVolume(ctx context.Context, req *quai.NFSPersistentVolumeReq, _ ...grpc.CallOption) (*quai.PersistentVolumeName, error) {
	pvReq := createNFSPVReq{
		Name: req.Name, Storage: req.Storage, Server: req.Server, Path: req.Server,
	}

	res, err := client.createNFSPersistentVolume(ctx, pvReq)
	if err != nil {
		return nil, err
	}

	pvRes := res.(createPVRes)
	return &quai.PersistentVolumeName{Value: pvRes.name}, pvRes.err
}

func (client *grpcClient) CreatePersistentVolumeClaim(ctx context.Context, req *quai.PersistentVolumeClaimReq, _ ...grpc.CallOption) (*quai.PersistentVolumeClaimName, error) {
	pvcReq := createPVCReq{
		Name: req.Name, Storage: req.Storage,
	}

	res, err := client.createNFSPersistentVolume(ctx, pvcReq)
	if err != nil {
		return nil, err
	}

	pvcRes := res.(createPVCRes)
	return &quai.PersistentVolumeClaimName{Value: pvcRes.name}, pvcRes.err
}

func (client *grpcClient) CreateDeployment(ctx context.Context, req *quai.DeploymentReq, _ ...grpc.CallOption) (*quai.DeploymentName, error) {
	resource := Resource{}
	if req.Resource != nil {
		resource.Memory = req.Resource.Memory
		resource.CPU = req.Resource.CPU
		resource.GPU = req.Resource.GPU
	}
	var volumes []*VolumeInfo
	for _, volume := range req.Volumes {
		volumes = append(volumes, &VolumeInfo{
			Name:      volume.Name,
			PVCName:   volume.PVCName,
			MountPath: volume.MountPath,
		})
	}
	deploymentReq := createDeploymentReq{
			Name:      req.Name,
			Replicas:  req.Replicas,
			Image:     req.Image,
			Resource:  &resource,
			Volumes:   volumes,
			Command:   req.Command,
			Arguments: req.Arguments,
	}

	res, err := client.createDeployment(ctx, deploymentReq)
	if err != nil {
		return nil, err
	}

	deploymentRes := res.(createDeploymentRes)
	return &quai.DeploymentName{Value: deploymentRes.name}, deploymentRes.err
}


func encodeCreateNFSPVRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(createNFSPVReq)
	return &quai.NFSPersistentVolumeReq{Name: req.Name, Storage: req.Storage, Server:req.Server, Path: req.Path}, nil
}

func encodeCreatePVCRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(createPVCReq)
	return &quai.PersistentVolumeClaimReq{Name: req.Name, Storage: req.Storage}, nil
}

func encodeCreateDeploymentRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(createDeploymentReq)
	resource := quai.Resource{}
	if req.Resource != nil {
		resource.Memory = req.Resource.Memory
		resource.CPU = req.Resource.CPU
		resource.GPU = req.Resource.GPU
	}
	var volumes []*quai.VolumeInfo
	for _, volume := range req.Volumes {
		volumes = append(volumes, &quai.VolumeInfo{
			Name:      volume.Name,
			PVCName:   volume.PVCName,
			MountPath: volume.MountPath,
		})
	}
	return &quai.DeploymentReq{
		Name:      req.Name,
		Replicas:  req.Replicas,
		Image:     req.Image,
		Resource:  &resource,
		Volumes:   volumes,
		Command:   req.Command,
		Arguments: req.Arguments,
	}, nil
}

func decodeCreateNFSPVResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*quai.PersistentVolumeName)
	return createPVRes{name: res.GetValue(), err: nil}, nil
}

func decodeCreatePVCResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*quai.PersistentVolumeClaimName)
	return createPVCRes{name: res.GetValue(), err: nil}, nil
}

func decodeCreateDeploymentResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*quai.DeploymentName)
	return createDeploymentRes{name: res.GetValue(), err: nil}, nil
}
