package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"stream/pb"
)

type HelloServiceImpl struct{}

func (p *HelloServiceImpl)Hello (
	ctx context.Context,args *pb.String)(*pb.String,error){

	fmt.Println("Hello receive:",args.GetValue())
 	reply := &pb.String{Value:"hello"+ args.GetValue()}
 	return reply,nil
}

//grpc stream
func (p *HelloServiceImpl)Channel(stream pb.HelloService_ChannelServer) error {
	for {
		args ,err := stream.Recv()
		if err != nil {
			if err==io.EOF{
				return nil
			}
			return err
		}

		fmt.Println("Channel receive:",args.GetValue())

		reply := &pb.String{Value:"stream hello:"+args.GetValue()}
		err = stream.Send(reply)
        if err !=nil {
        	return err
		}
	}
}

func main() {
	grpcServer := grpc.NewServer()

	pb.RegisterHelloServiceServer(grpcServer,new(HelloServiceImpl))

	lis,err := net.Listen("tcp",":1234")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer.Serve(lis)
}