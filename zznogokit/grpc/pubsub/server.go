package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"pubsub/pb"
	"pubsub/service"
)

func main(){
	grpcServer := grpc.NewServer()
	pb.RegisterPubsubServiceServer(grpcServer,service.NewPubsubService())

	lis,err := net.Listen("tcp",":1234")
	if err!=nil{
		log.Fatal(err)
	}

	grpcServer.Serve(lis)
}
