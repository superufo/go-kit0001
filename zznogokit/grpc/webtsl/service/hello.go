package service

import (
	"context"
	"webtsl/pb"
)

type MyGrpcServer struct{}

func (s *MyGrpcServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}