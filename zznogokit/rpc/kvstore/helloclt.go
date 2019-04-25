package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/rpc"
	"time"
)

//type kv struct {
//	key string
//	value string
//}

func main() {
	client, err := rpc.Dial("tcp", "127.0.0.1:12344")
	defer client.Close()

	if err != nil {
		log.Fatal("dialing:", err)
	}

	for {
		var txt string
		txt = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

		len := len(txt)
		var (
			randStr      string
			randValueStr string
		)

		var txtByte [100]string
		for ky, val := range txt {
			//fmt.Println(ky,string(val))
			txtByte[ky] = string(val)
		}
		fmt.Println(txtByte)

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := 0; i < 6; i++ {
			if i == 0 {
				randStr = ""
				randValueStr = ""
			}

			randInt := r.Intn(int(len))
			randStr = randStr + txtByte[randInt]

			randvInt := r.Intn(int(len) - 3)
			randValueStr = randValueStr + txtByte[randvInt] + txtByte[randvInt+1] + txtByte[randvInt+2]
			if i==5 {
				fmt.Println(randValueStr)
			}
		}

		var reply string
		//kvsl := new(kv)
		//kvsl.key = randStr
		//kvsl.value = randValueStr

		err = client.Call("KVStoreService.Set", [2]string{randStr, randValueStr}, &reply)
		fmt.Println(reply)
		err = client.Call("KVStoreService.Get", randStr, &reply)

		//err = client.Call("KVStoreService.GetValue", randStr, &reply)
		//fmt.Printf("kvsl%+v",kvsl)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(reply)

		doClientWork(client, [2]string{randStr, randValueStr})

		time.Sleep(time.Second * 30)
	}
}

func doClientWork(client *rpc.Client, kvParameter [2]string) {
	go func() {
		var keyChange string
		err := client.Call("KVStoreService.Watch", 30, &keyChange)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("watch:", keyChange)
	}()

	fmt.Printf("kvParameter%+v", kvParameter)

	var reply string
	kvParameter[1] = kvParameter[1] + "-abc"
	err := client.Call("KVStoreService.Set", kvParameter, &reply)
	if err != nil {
		log.Fatal(err)
	}
}
