package main

import (
	. "fmt"
	"sync"
)

type Test struct {
	Koin [2]int
}

var  chana = make(chan Test,1)
//var  chanb = chan test{}

var wg sync.WaitGroup

func main(){
	go write()

	wg.Add(1)
	go test1111()

	var kk = <-chana
	Printf("%+v\n",kk)

	wg.Wait()
	Println("end \n")
}

func write(){
    Println("sssss \n")
	//var s2 test = test{21,55}
	var tex *Test =  new(Test)
	tex.Koin = [2]int{1,67}
	chana <- *tex
	//wg.Done()
	return
}

func test1111()  {
	//time.Sleep(2 * time.Second)
	Printf("exit test1111 \n")
	wg.Done()
}

