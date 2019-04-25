package main

import (
	"fmt"
	"kvstore/lib"
	"log"
	"net"
	"net/rpc"
	"time"
)

func main() {
	//注册为Rpc函数
	rpc.RegisterName("KVStoreService", new(lib.KVStoreService))

	listener, err := net.Listen("tcp", ":12344")
	if err != nil {
		log.Fatal("listener Tcp error :", err)
	}

	var i int32
	for {
		go func() {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatal("Accept error", err)
			}

			fmt.Println(i)
			i++
			//提供rpc服务
			rpc.ServeConn(conn)
		}()

		time.Sleep(time.Second*30)
	}
}
