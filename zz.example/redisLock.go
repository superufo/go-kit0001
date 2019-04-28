package main

import (
	"github.com/go-redis/redis"
	redis "github.com/kataras/iris/sessions/sessiondb/redis"
)

func incr(){
	 client:=  redis.NewClient(&redis.Options{
		 Addr: "localhost:6379",
		 Password: "", // no password set
		 DB: 0, // use default DB
	 })

	 var 


}