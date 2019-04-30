package main

import (
	"log"
	"context"
	"google.golang.org/grpc"

	"interceptor/pb"
)

var (
	port = ":5000"
)

func main(){


	conn,err := grpc.Dial("localhost"+port,grpc.WithInsecure())
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