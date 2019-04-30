package service

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	"auth/pb"
)

type MyGrpcServer struct{}

func (s *MyGrpcServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	md ,ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil,fmt.Errorf("missing credentials")
	}

	var (
		appid string
		appkey string
	)

	if val,ok := md["login"];ok{
		appid = val[0]
	}

	if val,ok := md["password"]; ok {
		appkey = val[0]
	}

	// redis databases
	if appid!="gopher" || appkey!="password" {
		fmt.Printf("invalid token: appid=%s,appkey=%s",appid,appkey)
		return nil,grpc.Errorf(codes.Unauthenticated,"invalid token: appid=%s,appkey=%s",appid,appkey)
	}

	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

type Authentication struct {
	Login string
	Password string
}

func (a *Authentication)GetRequestMetadata(context.Context,...string)(map[string]string,error){
	return map[string]string{"login":a.Login,"password":a.Password},nil
}

func (a *Authentication)RequireTransportSecurity()bool{
	return false
}