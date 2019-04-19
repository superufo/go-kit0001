package main

import (
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewMultipleHostsReverseProxy(targets []*url.URL) *httputil.ReverseProxy{
	director := func(req *http.Request){
		target := targets[rand.Int()%len(targets)]

		// Scheme 为 http https
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path
	}

	return &httputil.ReverseProxy{Director:director}
}


/**
  net/http/httputil/ReverseProxy 反向代理库
 */
func main(){
	//把反向代理发送到 go run proxy_ex.go
	proxy := NewMultipleHostsReverseProxy([]*url.URL{
		{
			Scheme: "http",
			Host:   "127.0.0.1:5555",
		},
		{
			Scheme: "http",
			Host:   "127.0.0.1:8080",
		},
	})
	log.Fatal(http.ListenAndServe(":7775",proxy))
}