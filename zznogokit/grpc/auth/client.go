package main

import (
	"log"
	"context"
	"google.golang.org/grpc"

	"auth/pb"
	"auth/service"
)

var (
	port = ":5000"
)

func main(){
	auth := service.Authentication{
		Login :"gopher",
		Password:"password",
	}

	conn,err := grpc.Dial("localhost"+port,grpc.WithInsecure(),grpc.WithPerRPCCredentials(&auth))
    if err !=nil {
    	log.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)
    r,err := c.SayHello(context.Background(),&pb.HelloRequest{Name:"gopher"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("doClientWork: %s", r.Message)
}