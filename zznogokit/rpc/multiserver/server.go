package main

import (
	"log"
	"net"
	"net/rpc"

	"multiserver/services"
)

func main(){
	listener ,err := net.Listen("tcp",":20001")

	if err != nil {
		log.Fatal(" ListenTCP error: ",err)
	}

	// 为每一个 网络客户端  分配一个rpcServer 服务端
	for {
		conn,err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error",err)
		}

		go func() {
			defer conn.Close()

			p := rpc.NewServer()
			p.Register(&services.HelloService{
				conn:conn,
				isLogin:true,
			})

			rpc.ServeConn(conn)
		}()
	}

}

