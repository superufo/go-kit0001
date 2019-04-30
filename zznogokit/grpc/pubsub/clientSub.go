package main

import (
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"pubsub/pb"

	"context"
)

func main(){
	conn,err := grpc.Dial("127.0.0.1:1234",grpc.WithInsecure())
	if err !=nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewPubsubServiceClient(conn)
	stream,err := client.Subscribe(context.Background(),&pb.String{Value:"golang:"})
	if err!=nil {
		log.Fatal(err)
	}

	for{
		reply ,err := stream.Recv()
		if err!=nil {
			if err==io.EOF{
				break
			}
			log.Fatal(err)
		}

		fmt.Println(reply.GetValue())
	}
}

