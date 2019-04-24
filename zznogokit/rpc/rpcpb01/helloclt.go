package main

import (
	"fmt"
	"log"
	"net/rpc"
)

func main(){
	client ,err := rpc.Dial("tcp","127.0.0.1:12344")
	if err!=nil {
		log.Fatal("dialing:",err)
	}

	var reply string
	err = client.Call("HelloService.Hello","hello",&reply)
	if err !=nil {
		log.Fatal(err)
	}

	fmt.Println(reply)
}


