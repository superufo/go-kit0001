package example01

import (
	"context"
	"example01/pb"
	"log"
	"net"
	"google.golang.org/grpc"
)

type HelloServiceImpl struct{}

func (p *HelloServiceImpl) Hello(
	ctx context.Context, args *pb.String,
) (*pb.String, error) {
	reply := &pb.String{Value: "hello:" + args.GetValue()}
	return reply, nil
}

func main(){
	grpcServer := grpc.NewServer()
	pb.RegisterHelloServiceServer(grpcServer,new(HelloServiceImpl))

	lis,err := net.Listen("tcp",":1234")
	if err!=nil {
		log.Fatal(err)
	}

	grpcServer.Serve(lis)
}