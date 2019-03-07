package grpc

import (
	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/hykuan/k8s-client-example"
	"github.com/hykuan/k8s-client-example/models"
)

var _ quai.ModelServiceClient = (*grpcClient)(nil)

type grpcClient struct {
	startTraining endpoint.Endpoint
}

// NewClient returns new gRPC client instance.
func NewClient(conn *grpc.ClientConn) quai.ModelServiceClient {
	svcName := "quai.ModelService"

	return &grpcClient{
		startTraining: kitgrpc.NewClient(
			conn,
			svcName,
			"StartTraining",
			encodeStartTrainingRequest,
			decodeStartTrainingResponse,
			quai.PersistentVolumeName{},
		).Endpoint(),
	}
}

func (client *grpcClient) StartTraining(ctx context.Context, req *quai.TrainingReq, _ ...grpc.CallOption) (*quai.Training, error) {
	trainingReq := trainingReq{
		training: models.Training{
			Name:  req.Name,
			Image: req.Image,
			DataSet: &models.MountedPersistentVolumeClaim{
				PVCName:   req.DataSet.PVCName,
				MountPath: req.DataSet.MountPath,
			},
			Model: &models.MountedPersistentVolumeClaim{
				PVCName:   req.Model.PVCName,
				MountPath: req.Model.MountPath,
			},
			GPU:       req.GPU,
			Command:   req.Command,
			Arguments: req.Arguments,
		},
	}

	res, err := client.startTraining(ctx, trainingReq)
	if err != nil {
		return nil, err
	}

	trainingRes := res.(trainingRes)
	return &quai.Training{Value: trainingRes.name}, trainingRes.err
}

func encodeStartTrainingRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(trainingReq)
	return &quai.TrainingReq{
		Name:  req.training.Name,
		Image: req.training.Image,
		DataSet: &quai.MountedPersistentVolumeClaim{
			PVCName:   req.training.DataSet.PVCName,
			MountPath: req.training.DataSet.MountPath,
		},
		Model: &quai.MountedPersistentVolumeClaim{
			PVCName:   req.training.Model.PVCName,
			MountPath: req.training.Model.MountPath,
		},
		GPU:       req.training.GPU,
		Command:   req.training.Command,
		Arguments: req.training.Arguments,
	}, nil
}

func decodeStartTrainingResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*quai.Training)
	return trainingRes{name: res.GetValue(), err: nil}, nil
}
