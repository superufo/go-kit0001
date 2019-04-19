package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

type Pxy struct {}

/***
标准库
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
*/
func (p *Pxy) ServeHTTP (rw http.ResponseWriter, req *http.Request){
		fmt.Printf("Received request %s %s %s \n",req.Method,req.Host,req.RemoteAddr)
	    transport := http.DefaultTransport

		//step 1
	    outReq := new(http.Request)
	    *outReq = *req

	    /***每个代理服务器会在 X-Forwarded-For 头部填上前一个节点的 ip 地址，
	    这个地址可以通过 TCP 请求的 remote address 获取**/
	    /** X-Forwarded-For: client, proxy1, proxy2 **/
	    if clientIP,_, err := net.SplitHostPort(req.RemoteAddr);err == nil{
	    	if prior,ok := outReq.Header["X-Forward-For"]; ok {
	    		clientIP = strings.Join(prior,",")+ ","+clientIP
			}

	    	outReq.Header.Set("X-Forwarded-For",clientIP)
		}

	    //step2
	    res,err := transport.RoundTrip(outReq)
	    if err != nil {
	    	rw.WriteHeader(http.StatusBadGateway)
		}

	    //step 3
	    for key,value := range res.Header {
			for _, v := range value {
				rw.Header().Add(key,v)
			}
		}

	    rw.WriteHeader(res.StatusCode)
	    io.Copy(rw,res.Body)
	    res.Body.Close()
}

func main(){
	//var  pxy = new(Pxy)
	fmt.Println("server on 5555")
	http.Handle("/", &Pxy{})
	http.ListenAndServe("0.0.0.0:5555", nil)
}