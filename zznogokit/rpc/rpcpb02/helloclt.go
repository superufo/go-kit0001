package main

import (
	"fmt"
	"log"
	"math/rand"
	"rpcpb02/servers"
	"time"
)

func main(){
	client ,err := servers.DialHelloService("tcp","127.0.0.1:12344")
	if err!=nil {
		log.Fatal("dialing:",err)
	}

	for {
		var txt string
		txt = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

		//txtByte :=  str2bytes(txt)
		len := len(txt)
		var randStr string
		//var buffer bytes.Buffer
		var txtByte = []rune(txt)
		for i:=0;i<15;i++{
			randInt := rand.Int31n(int32(len))
			fmt.Println(txtByte[randInt]+"\n")
			randStr = randStr + string(txtByte[randInt])
		}

		var reply string
		err = client.Hello(randStr, &reply)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(reply)

		time.Sleep(time.Second*30)
	}
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


