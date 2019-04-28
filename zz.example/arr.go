package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
)

var msg string
var done = make(chan bool)

func main() {
	runtime.GOMAXPROCS(1)  //应用程序何以在运行期间设置运行时系统中得P最大数量

	//当参数的可变参数是空接口类型时，传人空接口的切片时需要注意参数展开的问题
	var a = []interface{}{1, 2, 3}
	fmt.Println(a)
	fmt.Println(a...)

	//在函数调用参数中，数组是值传递，无法通过修改数组类型的参数返回结果
	x := [3]int{1, 2, 3}
	func(arr [3]int) {
		arr[0] = 7
		fmt.Println(arr)
	}(x)
	fmt.Println(x)

	//map是一种hash表实现，每次遍历的顺序都可能不一
	m := map[string]string{
		"1": "1",
		"2": "2",
		"3": "3",
	}
	for k, v := range m {
		println(k, v)
	}

	//不同Goroutine之间不满足顺序一致性内存模型 要实现打印使用信号阻塞
	go setup()
	<-done
	fmt.Println(msg)

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(i)
		}
	}()

	//切片会导致整个底层数组被锁定，底层数组无法释放内存。如果底层数组较大会对内存产生很大的压力。
	headerMap := make(map[string][]byte)
	for i := 0; i < 5; i++ {
		name := "D:\\test.txt"
		data, err := ioutil.ReadFile(name)
		if err != nil {
			log.Fatal(err)
		}
		//如下  data内存 将无法释放
		//headerMap[name] = data[:1]
		//解决的方法是将结果克隆一份，这样可以释放底层的数组
		headerMap[name] = append([]byte{}, data[:1]...)  //append用来将元素添加到切片末尾并返回结果
		log.Println("headerMap:%+V",headerMap) //参数+V 打印接口
	}

	//recover捕获的是祖父级调用时的异常，直接调用时无效：必须在defer函数中直接调用才有效：
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("捕获到的错误：%s\n", r)
		}
	}()
	panic("1")

	//defer在函数退出时才能执行，在for执行defer会导致资源延迟释放   注释上面panic代码 执行下面代码
	//错误的代码
	//for i := 0; i < 5; i++ {
	//	f, err := os.Open("/path/to/file")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	defer f.Close()
	//}
	//解决的方法可以在for中构造一个局部函数，在局部函数内部执行defer：
	for i := 0; i<5; i++ {
		func() {
			f, err := os.Open("./text.txt")
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
		}()
	}

	//runtime.Goexit函数被调用后，会立即使调用他的Groution的运行被终止，但其他Goroutine并不会受到影响
	//runtime.Gosched函数的作用是暂停调用他的Goroutine的运行，调用他的Goroutine会被重新置于Gorunnable状态，并被放入调度器可运行G队列中
	//此处 若注释 runtime.Gosched()cpu 将会为for独占Goroutine 将不会运行    Goroutine是协作式抢占调度，Goroutine本身不会主动放弃CPU
	for {
		runtime.Gosched()
	}
	//此处 后面的代码不会执行

	//是通过阻塞的方式避免CPU占用
	//select{}
}

func setup() {
	msg = "hello, world"
	done <- true
}

/***
Go语言是带内存自动回收的特性，因此内存一般不会泄漏。但是Goroutine确存在
泄漏的情况，同时泄漏的Goroutine引用的内存同样无法被回收
上面的程序中后台Goroutine向管道输入自然数序列，main函数中输出序列。但是
当break跳出for循环的时候，后台Goroutine就处于无法被回收的状态了
 */
//func error_leak(){
//	ch := func() <-chan int {
//		ch := make(chan int)
//		go func() {
//			for i := 0; ; i++ {
//				ch <- i
//			}
//		} ()
//		return ch
//	}()
//
//	for v := range ch {
//		fmt.Println(v)
//		if v == 5 {
//			break
//		}
//	}
//}

/**
可以通过context包来避免这个问题（Goroutine内存泄漏的情况)
当main函数在break跳出循环时，通过调用 cancel() 来通知后台Goroutine退出，
这样就避免了Goroutine的泄漏
 */
 //func ok_leak(context) <-chan int {
	//ctx, cancel := context.WithCancel(context.Background())
	//ch := func(ctx context.Context) <-chan int {
	//	 ch := make(chan int)
	//	 go func() {
	//		 for i := 0; ; i++ {
	//			 select {
	//			 case <- ctx.Done():
	//				 return
	//			 case ch <- i:
	//			 }
	//		 }
	//	 } ()
	//	 return ch
	//}(ctx)
 //
	//for v := range ch {
	//	 fmt.Println(v)
	//	 if v == 5 {
	//		 cancel()
	//		 break
	//	 }
	//}
 //}


