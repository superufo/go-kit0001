package main

import (
	"log"
	"net"
	"net/rpc"
	. "./services/hello"
)
func main(){
	rpc.RegisterName("HelloService",new(HelloService))
    listener ,err := net.Listen("tcp",":1234")
    if err != nil {
    	log.Fatal("listen tcp error : ",err)
	}

    conn,err := listener.Accept()
    if err!=nil{
		log.Fatal("listen accept error : ",err)
	}

    rpc.ServeConn(conn)
}