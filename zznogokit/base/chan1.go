package main

import "fmt"

func main() {
	done := make(chan int, 10) // 带 10 个缓存
	// 开N个后台打印线程
	for i := 0; i < cap(done); i++ {
		go func(){
			fmt.Println("你好, 世界")
			done <- 1
		}()
	}
	// 等待N个后台线程完成
	for i := 0; i < cap(done); i++ {
		<-done
	}
}
/*
 wg.Add(1) 用于增加等待事件的个数，必须确保在后台线程启动之前执行
（如果放到后台线程之中执行则不能保证被正常执行到） 。当后台线程完成打印工
作之后，调用 wg.Done() 表示完成一个事件。 main 函数的 wg.Wait() 是等待
全部的事件完成
 */