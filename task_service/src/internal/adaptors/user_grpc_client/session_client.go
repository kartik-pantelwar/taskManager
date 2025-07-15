package client

import (
	pb "task_service/src/internal/interfaces/input/grpc/generated/generated"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewSessionValidatorClient(addr string) (pb.SessionValidatorClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return pb.NewSessionValidatorClient(conn), nil
}
