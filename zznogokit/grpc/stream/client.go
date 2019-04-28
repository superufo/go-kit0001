package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"math/rand"
	"strconv"
	"time"

	"stream/pb"
)

func main(){
	conn,err := grpc.Dial("127.0.0.1:1234",grpc.WithInsecure())
	if err!=nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewHelloServiceClient(conn)
	reply,err := client.Hello(context.Background(),&pb.String{Value:"hello"})
	if err!=nil {
		log.Fatal(err)
	}
	fmt.Println(reply.GetValue())

	//grpc stream
	stream,err := client.Channel(context.Background())
	if err!=nil {
		log.Fatal(err)
	}

	go func(){
		for {
			rand.New(rand.NewSource(time.Now().UnixNano()))
			rint := rand.Int63()

			//fmt.Println(rint)
			if err := stream.Send(&pb.String{Value:strconv.FormatInt(rint,10)});err!=nil{
				log.Fatal(err)
			}
			time.Sleep(time.Second*5)
		}
	}()

	for{
		reply,err := stream.Recv()
		if err !=nil {
			if err == io.EOF{
				break
			}
		}
		fmt.Println(reply.GetValue())
	}



}

