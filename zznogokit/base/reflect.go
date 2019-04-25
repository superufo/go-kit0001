package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main(){
	 var ag = [...]string{"test","go","tiandi","健康"}

	 fmt.Println("len(s):",(*reflect.StringHeader)(unsafe.Pointer(&ag)).Len)
	 fmt.Printf("b: %T\n", ag) // b: [3]int
	 fmt.Printf("b: %#v\n", ag)

	for i := range ag {
		fmt.Printf("a[%d]: %d\n", i, ag[i])
	}

	//z指针转化
	//任何指针都可以转换为unsafe.Pointer
	//unsafe.Pointer可以转换为任何指针
	//uintptr可以转换为unsafe.Pointer
	//unsafe.Pointer可以转换为uintptr
	i:= 10
	ip:=&i
	var fp *float64 = (*float64)(unsafe.Pointer(ip))
	*fp = *fp * 3
	fmt.Println(i)
}
