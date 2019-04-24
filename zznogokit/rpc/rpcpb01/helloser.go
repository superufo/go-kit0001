package main

import (
	"log"
	"net"
	"net/rpc"
	"rpcpb01/servers"
)

func main() {
	//注册为Rpc函数
	rpc.RegisterName("HelloService",new(servers.HelloService))

	listener,err := net.Listen("tcp",":12344")
	if err !=nil {
		log.Fatal("listener Tcp error :",err)
	}

	conn ,err := listener.Accept()
	if err!=nil {
		log.Fatal("Accept error",err)
	}

	//提供rpc服务
	rpc.ServeConn(conn)
}