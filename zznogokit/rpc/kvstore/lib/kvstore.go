package lib

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type KVStoreService struct {
	m map[string]string
	filter map[string]func(key string)
	mu sync.Mutex
}

func NewKVStoreService() *KVStoreService{
	return &KVStoreService{
		m: make(map[string]string),
		filter :make(map[string]func(key string)),
	}
}

func (p *KVStoreService) Get(key string ,value *string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	//获取key 的value
	if v,ok := p.m[key];ok {
		*value = v
		return nil
	}

	return fmt.Errorf("not found")
}

//func (p *KVStoreService) GetValue(key string)(val string) {
//	p.mu.Lock()
//	defer p.mu.Unlock()
//
//	//获取key 的value
//	if v, ok := p.m[key]; ok {
//		val = v
//		//err = nil
//	}else {
//		val = ""
//		//err = errors.New("not found")
//	}
//	return
//}

/***
kv  key 和 value
*/
func (p *KVStoreService) Set (kv [2]string,reply *string) error {
	fmt.Printf("kvsl%+v",kv)
	p.mu.Lock()
	defer p.mu.Unlock()

	key,value := kv[0],kv[1]
	if oldValue := p.m[key]; oldValue!= value {
		for _,fn := range p.filter{
			//调用filter中函数
			fn(key)
		}
	}

	*reply = "ok"
	p.m[key] = value
	return nil
}

func (p *KVStoreService) Watch (timeoutSecond int, keyChanged *string)error {
	id := fmt.Sprintf("watch-%s-%03d",time.Now(),rand.Int())
	ch := make(chan string ,10)

	p.mu.Lock()
	p.filter[id] = func(key string){ ch <-key }
	p.mu.Unlock()

	select {
		case <-time.After(time.Duration(timeoutSecond)*time.Second):
			return fmt.Errorf("timeout")
		//当45行调用的时候这里会执行
		case key:=<-ch:
			*keyChanged =key
			return nil
	}

	return nil
}