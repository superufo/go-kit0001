package main

import (
	"fmt"
	"time"
)

func main(){
	var  fillInterval = time.Microsecond * 10
	var  capacity =100
	var tokenBucket = make (chan struct{},capacity)

	//time.NewTicker 定时器
	fillToken := func(){
		//NewTimer 创建一个 Timer，它会在最少过去时间段 d 后到期，向其自身的 C 字段发送当时的时间
		ticker := time.NewTicker(fillInterval)
		for {
			select {
				case <- ticker.C:
					select{
						case tokenBucket <- struct{}{}:
						default :
					}
					fmt.Println(" current token cnt:",len(tokenBucket),time.Now())
			}
		}
	}

	go fillToken()
	time.Sleep(time.Hour)
}

func TakeAvailable(block bool) bool{
	var takenResult bool
	if block {
		select {
		case <-tokenBucket:
			takenResult = true
		}
	} else {
		select {
		case <-tokenBucket:
			takenResult = true
		default:
			takenResult = false
		}
	}
	return takenResult
}


