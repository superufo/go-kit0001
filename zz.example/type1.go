package main

import "fmt"

type Status []byte

func main(){
	kk := Status{1,2,3,5}
	lk := ([]byte)(kk)

	测试:= []string{"xcxc","cvcv石狮市"}

	for _,一个字符串:= range 测试 {
		fmt.Printf("一个字符串:%s\n",一个字符串)
	}

	for _,c:= range lk {
		fmt.Printf("c:%d\n",c)
	}

	fmt.Printf("KK:%s",lk)
	var m=string(kk)
	fmt.Printf("KK:%s",m)
}

