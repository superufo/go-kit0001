package main

import (
	. "fmt"
	"sync"
	"time"
)

var a=[...]int{1,6:2,9}

func main(){
	b := &a
	Printf("a:%+v\n",a)

	var wg  sync.WaitGroup

	for  k,v := range b {
		Printf("k,v : %d   %d \n",k,v)
	}

	var times [5][0]int
	for range times {
		Println("hello")
	}

	var ms chan int
	go func() {
		wg.Add(1)
		Printf("gproutine start....ms.....\n")
		ms <- 45
	}()

	// 管道数组
	//var chanList = [2]chan int{}
    //Printf("len(chanList) %d\n",len(chanList))
	var  chanList = make(chan int,1)
    go func() {
		Printf("\ngproutine start")
		c := time.Tick(1 * time.Second)
		for i:=20;i>0;i--  {
			wg.Add(1)
			Printf("chan: %s", "goroutine.chanList..\n")
			Printf("chan:%d \n", i)
			//var kk = [2]int{1,2}
			//chanList <- kk
			chanList <- i
			<- c
		}
	}()

	//for  h := range  chanList {
	//	Printf("chan h---:%+v", h)
	//}
	 var aaa = <-chanList
	Printf("test2 ------  %d \n",aaa)

	var aaag = <-ms
	Printf("test3 ------  %d \n",aaag)

	for {
		select {
			case <-ms:
				wg.Done()
				Printf("test2 \n")
				break
			case <-chanList:
				wg.Done()
				Printf("test \n")
				break
			default:
		}
	}

	wg.Wait()
	Printf("main end")
}


