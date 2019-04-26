package main

import (
	"fmt"
	"log"
	"net/rpc"
)

func main(){
	client ,err := rpc.Dial("tcp","127.0.0.1:20001")
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

// var a = []byte("hello world ")  var b = string(a) 直接转换 是通过 数据copy
// 在并发达到 百万级别性能不高
// gdb 调试知道 string 看作 [2]uintptr []byte 则是[3]unitptr
//func str2bytes(s string) []byte{
//     x := (*[2]uintptr)(*unsafe.Pointer(&s))
//     h :=  [3]uintptr{x[0],x[1],x[1]}
//	 return *(*[]byte)(unsafe.Pointer(&h))
//}
//
//func byte2string(b []byte) string{
//	return *(*string)(unsafe.Alignof(&b))
//}


