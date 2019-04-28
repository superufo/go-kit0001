package main

import (
	"sync"
)

// 全局变量

var counter int
var l = sync.Mutex{}
func main() {
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.Lock()
			counter++
			println("current counter", counter)
			l.Unlock()
		}()
	}

	wg.Wait()
	println(counter)
}
