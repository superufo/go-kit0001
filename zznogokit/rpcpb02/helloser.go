package main

import (
	"log"
	"net"
	"net/rpc"
	"rpcpb02/servers"
)
type HelloService struct {}

func (P *HelloService) Hello (request string,reply *string) error {
	*reply =  "hello world, request:" +  request
	return nil
}

func main() {
	//注册为Rpc函数
	servers.RegistertHelloService(new(HelloService))

	listener,err := net.Listen("tcp",":12344")
	if err !=nil {
		log.Fatal("listener Tcp error :",err)
	}

   //支持多tcp连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error", err)
		}

		//提供rpc服务
		rpc.ServeConn(conn)
	}
}