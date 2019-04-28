package main

import (
	"google.golang.org/grpc"
	"log"
	"fmt"
	"context"

	"hello/pb"
)

func main(){
	conn,err := grpc.Dial("127.0.0.1:1234",grpc.WithInsecure())
	if err !=nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := pb.NewHelloServiceClient(conn)
	reply,err := client.Hello(context.Background(),&pb.String{Value:"hello"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(reply.GetValue())
}