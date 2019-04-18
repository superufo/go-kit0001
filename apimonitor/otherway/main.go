package main

import (
	"github.com/rcrowley/go-metrics"
	"time"
	"os"
	"log"
)

func main(){
	g := metrics.NewGauge()
	metrics.Register("bar", g)
	g.Update(1)


	go metrics.Log(metrics.DefaultRegistry,
		1 * time.Second,
		log.New(os.Stdout, "metrics: ", log.Lmicroseconds))


	var j int64
	j = 1
	for true {
		time.Sleep(time.Second * 1)
		g.Update(j)
		j++
	}
}