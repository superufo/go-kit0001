package main

import (
	"fmt"
	"sync"
)

type Lock struct {
	c chan struct{}
}

func NewLock() Lock {
	var l Lock
    l.c= make(chan struct{},1)
    l.c<-struct {}{}
    return  l
}

func (l Lock) Lock() bool{
	lockResult := false
	select {
		case  <-l.c :
			lockResult=true
	    default:
	}
	return lockResult
}

func (l Lock) Unlock() {
	l.c <- struct {}{}
}

var count int
func main(){
	var l =NewLock()
    var wg  sync.WaitGroup

    for i:=0 ;i<10 ;i++ {
    	wg.Add(1)
    	go func() {
    		defer wg.Done()

    		if !l.Lock(){
    			println("获取锁失败\n")
    			return
			}
			count++
			fmt.Printf("count:%d\n",i)
    		l.Unlock()
		}()
	}

    wg.Wait()
}



