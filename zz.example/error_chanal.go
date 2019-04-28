package main 

import (
   "fmt"
   "time"
)


var done = make(chan bool)
var msg string
func aGoroutine() {
	time.Sleep(10000)
	msg = "hello, world" 
	<-done
} 

func main() {
	go aGoroutine()
	done <- true
	time.Sleep(10000)
	fmt.Println(msg)
}
