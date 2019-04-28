package main

import (
	"fmt"
)

/**
通道变量 会阻塞
当 <-done 执行时，必然要求 done <- 1 也已经执行。根据同一个Gorouine依然
满足顺序一致性规则，我们可以判断当 done <- 1 执行时， println("你好, 世
界") 语句必然已经执行完成了。因此，现在的程序确保可以正常打印结果。
 */
func main() {
	done := make(chan int)
	go func(){
		fmt.Println("你好, 世界")
		done <- 1
	}()

	<-done
}

//func main() {
//	var mu sync.Mutex
//	mu.Lock()
//	go func(){
//		println("你好, 世界")
//		mu.Unlock()
//	}()
//	mu.Lock()
//}
//可以确定后台线程的 mu.Unlock() 必然在 println("你好, 世界") 完成后发生
//（同一个线程满足顺序一致性） ， main 函数的第二个 mu.Lock() 必然在后台线
//程的 mu.Unlock() 之后发生 ( ，当第二次加锁时会因为锁已经被占用（不是递归锁） 而阻塞)（ sync.Mutex 保证） ，此时后台线程的打印工作已
//经顺利完成了。


