package main

import (
	"fmt"
	"jsonrpc/servers"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main() {
	//注册为Rpc函数
	rpc.RegisterName("HelloService", new(servers.HelloService))

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listener Tcp error :", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error", err)
		}

		fmt.Println("RemoteAddr: ",conn.RemoteAddr(),"LocalAddr: ",conn.LocalAddr())
		//提供rpc服务  func ServeConn(conn io.ReadWriteCloser)
		go 	rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
