package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)


/**
type serverRequest struct {
	Method string           `json:"method"`
	Params *json.RawMessage `json:"params"`
	Id     *json.RawMessage `json:"id"`
}

type serverResponse struct {
	Id     *json.RawMessage `json:"id"`
	Result interface{}      `json:"result"`
	Error  interface{}      `json:"error"`
}
 */
func main(){
	conn ,err := net.Dial("tcp","127.0.0.1:1234")
	if err!=nil {
		log.Fatal("dialing:",err)
	}
	//NewClientCodec(conn io.ReadWriteCloser) rpc.ClientCodec {
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))

	var reply string
	err = client.Call("HelloService.Hello","hello",&reply)
	if err !=nil {
		log.Fatal(err)
	}

	fmt.Println(reply)
}


