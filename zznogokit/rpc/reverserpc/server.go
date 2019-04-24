package main

import (
	"net"
	"net/rpc"
	"reverserrpc/service/HelloService"
	"time"
)

//服务器端  内网主动拨外网建立连接
func main(){
	//rpc 服务器端
	rpc.Register(new(HelloService))

	for{
		//tcp客户端
		conn,_:= net.Dial("tcp","47.112.111.171:1234")

		if conn==nil {
			time.Sleep(time.Second)
			continue
		}

		//rpc 服务器端
		rpc.ServeConn(conn)
		conn.Close()
	}
}

