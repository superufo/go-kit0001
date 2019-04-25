package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Producer(factor int,out chan<- int){
	for i:=0;;i++{
		out <- i*factor
	}
}

func Consumer(in <-chan int){
	for v := range in {
		fmt.Println(v)
	}
}

func main(){
	ch:=make(chan int,64)

	go Producer(3,ch)
	go Producer(5,ch)
	go Consumer(ch)

	//运行一定时间后退出
	time.Sleep(5*time.Second)

	// Ctrl+C 退出
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("quit (%v)\n", <-sig)

	////合建chan
	//c := make(chan os.Signal)
	////监听指定信号 ctrl+c kill
	////signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
	//signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	////阻塞直到有信号传入
	//fmt.Println("启动")
	////阻塞直至有信号传入
	//s := <-c
	//fmt.Println("退出信号", s)
}