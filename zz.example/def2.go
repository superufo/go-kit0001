package main

import (
	"fmt"
	"time"
	"os"
)

func launch() {
	fmt.Println("nuclear launch detected")
}

func commencingCountDown(canLunch chan int) {
	c := time.Tick(1 * time.Second)
	fmt.Printf("%+v\n",c)
	for countDown := 20; countDown > 0; countDown-- {
		fmt.Println(countDown)
		<- c
	}
	canLunch <- -1
}

func isAbort(abort chan int) {
	os.Stdin.Read(make([]byte, 1))
	abort <- -1
}

func main() {
	fmt.Println("Commencing coutdown")

	abort := make(chan int)
	canLunch := make(chan int)
	go isAbort(abort)
	go commencingCountDown(canLunch)

	 //for _,tmp :=  range <-canLunch {
		// fmt.Printf("recv:%+v\n", tmp)
	 //}

	//for {
	//	i, ok :=  <-canLunch
	//	if !ok {
	//		break
	//	}
	//	fmt.Println(i)
	//}

	select {
		case <- canLunch:

		case <- abort:
			fmt.Println("Launch aborted!")
			return
	}
	launch()
}