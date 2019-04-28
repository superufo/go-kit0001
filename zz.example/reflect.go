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

	BytoToString([]byte{'1','f','h'})


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


	s := struct {
		a byte
		b byte
		c byte
		d int64
	}{0, 0, 0, 0}

	// 将结构体指针转换为通用指针
	p := unsafe.Pointer(&s)
	// 保存结构体的地址备用（偏移量为 0）
	up0 := uintptr(p)
	// 将通用指针转换为 byte 型指针
	pb := (*byte)(p)
	// 给转换后的指针赋值
	*pb = 10
	// 结构体内容跟着改变
	fmt.Println(s)

	// 偏移到第 2 个字段
	up := up0 + unsafe.Offsetof(s.b)
	// 将偏移后的地址转换为通用指针
	p = unsafe.Pointer(up)
	// 将通用指针转换为 byte 型指针
	pb = (*byte)(p)
	// 给转换后的指针赋值
	*pb = 20
	// 结构体内容跟着改变
	fmt.Println(s)

	// 偏移到第 3 个字段
	up = up0 + unsafe.Offsetof(s.c)
	// 将偏移后的地址转换为通用指针
	p = unsafe.Pointer(up)
	// 将通用指针转换为 byte 型指针
	pb = (*byte)(p)
	// 给转换后的指针赋值
	*pb = 30
	// 结构体内容跟着改变
	fmt.Println(s)

	// 偏移到第 4 个字段
	up = up0 + unsafe.Offsetof(s.d)
	// 将偏移后的地址转换为通用指针
	p = unsafe.Pointer(up)
	// 将通用指针转换为 int64 型指针
	pi := (*int64)(p)
	// 给转换后的指针赋值
	*pi = 40
	// 结构体内容跟着改变
	fmt.Println(s)
}

func BytoToString(barr []byte) (res string) {
	data := make([]byte,len(barr))

	//for  i:=0; i<=len(barr);i++  {
	//	data[i] = barr[i]
	//}
	for i,k :=  range barr {
		data[i] = k
	}

	hdr := (*reflect.StringHeader)(unsafe.Pointer(&res))
    hdr.Data = uintptr(unsafe.Pointer(&data[0]))  //unsafe.Pointer类似 C 语言中的 void* 万能指针，能安全持有对象
	hdr.Len = len(barr)

	fmt.Println(hdr)
	return res
}
