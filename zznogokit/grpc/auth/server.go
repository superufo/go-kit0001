package main

import (
	"auth/pb"
	"auth/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	port = ":5000"
)

func main(){
	server := grpc.NewServer()
	pb.RegisterGreeterServer(server,new(service.MyGrpcServer))

	lis ,err := net.Listen("tcp",port)
	if err !=nil {
		log.Panicf("could not list on %s: %s", port, err)
	}

	if err := server.Serve(lis); err != nil {
		log.Panicf("grpc serve error: %s", err)
	}
}