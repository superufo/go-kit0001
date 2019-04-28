package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"reflect"
	"unsafe"
)


type ByteStruct  struct{
	RD []byte
}

var ptr chan ByteStruct
var sw = make(chan bool)

func main() {
	// 字符串数组
	var s1 = [2]string{"hello", "world"}
	var s2 = [...]string{"你好", "世界"}
	var s3 = [...]string{1: "世界", 0: "你好", }
	// 结构体数组
	var line1 [2]image.Point
	var line2 = [...]image.Point{image.Point{X: 0, Y: 0}, image.Point{X: 1, Y: 1}}
	var line3 = [...]image.Point{{0, 0}, {1, 1}}
	// 图像解码器数组
	var decoder1 [2]func(io.Reader) (image.Image, error)
	var decoder2 = [...]func(io.Reader) (image.Image, error){
		png.Decode,
		jpeg.Decode,
	}
	// 接口数组
	var unknown1 [2]interface{}
	var unknown2 = [...]interface{}{123, "你好"}


	//var d [0]int // 定义一个长度为0的数组
	//var e = [0]int{} // 定义一个长度为0的数组
	//var f = [...]int{} // 定义一个长度为0的数组
	//%T 或 %#v 谓词语法来打印数组的类型和详细信息
	fmt.Printf("s1: %T\n", s1)
	fmt.Printf("vars:%+v\n%+v\n%+v\n%+v\n%+v\n%+v\n%+v\n%+v\n%+v\n%+v\n\n",line1,line2,line3,decoder1,decoder2,unknown1,unknown2,s3,s1,s2)


	c2 := make(chan struct{})

	go func() {
		fmt.Println("c2")
		c2 <- struct{}{} // struct{}部分是类型, {}表示对应的结构体值
	}()

	<-c2

	s := "hello, world"
	fmt.Println("len(s):", (*reflect.StringHeader)(unsafe.Pointer(&s)).Len)

	go strtobyt("test00000")
	for {
		//var dd = <-ptr
		//fmt.Println( dd )
		fmt.Printf("ptr:::::%+v\n", <- ptr)
	}

	// 管道数组
	//var chanList = [2]chan int{}

	//<-chanList[0]
	//<-chanList[1]
	//go func(){
	//	fmt.Println("c3")
	//
	//	for i:=0; i<cap(chanList);i++{
	//		chanList[i]<-i
	//	}
	//}()
	//<-chanList[0]
	//<-chanList[1]

	//fmt.Printf("dd:%d \n", dd)

	<-sw
}


func strtobyt(str string)  {
	p := make([]byte, len(str))

	fmt.Println(str)
	for i:=0; i<len(str);i++{
		c := str[i]
		fmt.Printf("c:::::%v\n",c)
		p[i] = c
		//fmt.Printf("kk:::::%#d\n",i)
	}

	//var kkl = new(ByteStruct)
	//kkl.RD = p

	lp :=  ByteStruct{RD:p}

	fmt.Println(lp)
	fmt.Printf("p:::%+v\n",lp)
	fmt.Println(lp.RD)

	mp:= ByteStruct{RD:[]byte{1,2}}
	ptr <- mp

	sw <- true // 阻塞作用
}



