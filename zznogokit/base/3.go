package main

import (
	"sync"
	"sync/atomic"
	"fmt"
	"sync/mutex"
)

func main(){
	var s =Instance()
    s.Kt();
	fmt.Printf("\n %+v",s)
	fmt.Printf("\n %v",s)
}

type singleton struct {}
var (
	instance *singleton
	initialized uint32
	mu sync.Mutex
)

func  Instance() *singleton {
	if atomic.LoadUint32(&initialized) == 1 {
	  return instance
	}

	mu.Lock()
	defer mu.Unlock()
	if instance == nil {
		defer atomic.StoreUint32(&initialized, 1)
		instance = &singleton{}
	}
	return instance
}

func (*singleton) Kt() {
	var a = []int{1,2,3,4}

	for j,y := range a {
		fmt.Printf("\nj,y : %d-%d",j,y)
	}
}


