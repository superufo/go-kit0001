package main

import (
	"./impl/helloservice"
	"./pb/hellopro"
	"github.com/go-kit/kit/transport/grpc"
	"grpc/impl/helloservice"
	"grpc/pb/hellopro"
	"log"
	"net"
)

func main() {
	grpcServer := grpc.NewServer()

	hellopro.RegisterHelloServiceServer(grpcServer, new(helloservice.HelloServiceImpl))
	//RegisterHelloServiceServer(grpcServer, new(HelloServiceImpl))

	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer.Serve(lis)
}
