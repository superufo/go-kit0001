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

/**
服务端在循环中接收客户端发来的数据，如果遇到io.EOF表示客户端流被关闭，如
果函数退出表示服务端流关闭。生成返回的数据通过流发送给客户端，双向流数据
的发送和接收都是完全独立的行为。需要注意的是，发送和接收的操作并不需要一
一对应，用户可以根据真实场景进行组织代
 */
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

	//向服务端发送数据
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

	//然后在循环中接收服务端返回的数据
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

