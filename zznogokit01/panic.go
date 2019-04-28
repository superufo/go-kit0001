package main

import (
	"fmt"
	"log"
)

func main() {

	if r := recover(); r != nil {
		log.Fatal(r)
	}

	fmt.Printf("texx:%d",111)

	//panic(123)
	panic("system error")

	fmt.Printf("texx:%d",777)

	if r := recover(); r != nil {
		log.Fatal(r)
	}

}
