package service

import (
	"context"
	"interceptor/pb"
)

type MyGrpcServer struct{}

func (s *MyGrpcServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	 panic("debug")
	 return &pb.HelloReply{Message:"Hello: "+in.Name},nil
}