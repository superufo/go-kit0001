package main

import (
	"fmt"
	"log"
	"net/rpc"
)

func main(){
	client,err := rpc.Dial("tcp","127.0.0.7:1234")
	if err!=nil {
		log.Fatal("Erro :",err)
	}

	var reply string
	err = client.Call("HelloService.Hello","hello",&reply)
    if err!=nil {
		log.Fatal("Erro :",err)
	}

	fmt.Printf("reply :%s",reply)
	fmt.Println(reply)
}
