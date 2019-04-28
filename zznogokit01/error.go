package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

//Web框架一般会通过 recover 来防御性
//地捕获所有处理流程中可能产生的异常，然后将异常转为普通的错误返回。
func main(){
	//defer func() {
	//	if r := recover(); r != nil {
	//		fmt.Printf("error:%+v",r)
	//		log.Fatal(r)
	//	}
	//}()


	CopyFile("1bak.go","./1.log")
}

func CopyFile (dstName, srcName string) (written int64, err error) {
	fmt.Printf("error:%+v",srcName)


	src, err := os.Open(srcName)
	if err != nil {
		fmt.Println("11111111111")
		//return
	}

	defer src.Close()
	dst, err := os.Create(dstName)
	if err != nil {
		fmt.Println("2222222222")
        //return
	}

	defer func(){
		if r := recover(); r != nil {
			fmt.Printf("error:%+v",r)
			fmt.Errorf("JSON: internal error: %v", r)
			log.Fatal(r)
		}
		dst.Close()
	}()


	return io.Copy(dst, src)
}