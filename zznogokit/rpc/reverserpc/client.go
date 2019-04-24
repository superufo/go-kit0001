package main

import (
	"fmt"
	"github.com/prometheus/common/log"
	"net"
	"net/rpc"
)

//外网客户端
func main(){
	//tcp服务器
	listener ,err := net.Listen("tcp",":1234")
	if err != nil {
		log.Fatal("ListenTCP error",err)
	}

	clientChan := make(chan *rpc.Client)

	go func(){
		for{
			conn,err := listener.Accept()
			if err!=nil {
				log.Fatal("Accept error:",err)
			}

			//NewClient(conn io.ReadWriteCloser)
			//rpc客户端
			clientChan <- rpc.NewClient(conn)
		}
	}()

	doClientWork(clientChan)
}

func doClientWork(clientchan <-chan *rpc.Client){
	client <- clientchan
	defer client.Close()

	var reply string
	//rpc客户端
	err = client.call("HelloService.Hello","hello",&reply)

	if err !=nil{
		log.Fatal(err)
	}

	fmt.Println(reply)
}