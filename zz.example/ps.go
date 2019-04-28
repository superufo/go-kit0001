package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type  Goods struct{
	hight int
}

// strconv.Atoi() 字符转数字
func produces(fac int,out chan Goods){
	for i:=0;;i++ {
		var GoodsEx =  new(Goods)
		GoodsEx.hight = i*fac
		fmt.Println("生产 GoodsEx:",*GoodsEx)
		out <-  *GoodsEx
	}
}

func consumer(in chan Goods){
	for i:=0;;i++ {
		kk := new(Goods)
		*kk = <-in
		fmt.Println("消耗 i:",i,"in:",*kk)
	}
}

func main() {
	ch := make(chan Goods,64)
	go produces(3,ch)
	go produces(5,ch)
	go consumer(ch)

	// main 函数保存阻塞状态不退出，只有当用户输入 Ctrl-C 时才真正退出程序
	sig := make(chan os.Signal,1)
	signal.Notify(sig,syscall.SIGINT,syscall.SIGTERM)
	fmt.Printf("quit (%v)\n", <-sig)
}
