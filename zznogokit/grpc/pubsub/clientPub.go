package main

import (
	"google.golang.org/grpc"
	"log"
	"pubsub/pb"

	"context"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:1234", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := pb.NewPubsubServiceClient(conn)
	//实际工程中 这里是业务逻辑 然后发布
	_, err = client.Publish(context.Background(), &pb.String{Value: "golang:hello Go"})
	if err != nil {
		log.Fatal(err)
	}

	_,err = client.Publish(context.Background(),&pb.String{Value:"docker :hello Docker"})
	if err !=nil{
		log.Fatal(err)
	}
}
