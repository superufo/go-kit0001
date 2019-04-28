package main

import (
	"context"
	"log"
	"net"
	"google.golang.org/grpc"

	"hello/pb"
)

type HelloServiceImpl struct{}

func (p *HelloServiceImpl) Hello(
	ctx context.Context, args *pb.String,
) (*pb.String, error) {
	reply := &pb.String{Value: "hello world:" + args.GetValue()}
	return reply, nil
}

func main(){
	grpcServer := grpc.NewServer()
	pb.RegisterHelloServiceServer(grpcServer,new(HelloServiceImpl))

	lis,err := net.Listen("tcp",":1234")
	if err!=nil {
		log.Fatal(err)
	}

	grpcServer.Serve(lis)
}