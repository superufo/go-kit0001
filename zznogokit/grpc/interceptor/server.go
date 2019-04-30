package main

import (
	"fmt"
	"interceptor/pb"
	"interceptor/service"

	"google.golang.org/grpc"
	"log"
	"net"
	"context"
)

var (
	port = ":5000"
)

//要实现普通方法的截取器，需要为grpc.UnaryInterceptor的参数实现一个函数
//函数的ctx和req参数就是每个普通的RPC方法的前两个参数。第三个info参数表示
//当前是对应的那个gRPC方法，第四个handler参数对应当前的gRPC方法函数。上
//面的函数中首先是日志输出info参数，然后调用handler对应的gRPC方法函数
//截取器也非常适合前面对Token认证工
//gRPC框架中只能为每个服务设置一个截取器
//开源的grpc-ecosystem项目中的go-grpc-middleware包已经基于gRPC对截取器实现了链式截取器的支持
func filter(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handle grpc.UnaryHandler,
	)(resp interface{},err error){
		log.Println("filter: ",info)

		defer func(){
			if  r := recover();r!=nil {
				err = fmt.Errorf("panic :%v",r)
			}
		}()

		return handle(ctx,req)
}


func main(){
	//服务器在收到每个gRPC方法调用之前，会首先输出一行日志，然后再调用对方的方法
	server := grpc.NewServer(grpc.UnaryInterceptor(filter))
	pb.RegisterGreeterServer(server,new(service.MyGrpcServer))

	lis ,err := net.Listen("tcp",port)
	if err !=nil {
		log.Panicf("could not list on %s: %s", port, err)
	}

	if err := server.Serve(lis); err != nil {
		log.Panicf("grpc serve error: %s", err)
	}
}