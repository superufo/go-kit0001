package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/hashicorp/consul/api"
)

func main(){
	// 创建环境变量
	var (
		consulHost = flag.String("consul.host","47.112.111.171","")
		consulPort = flag.String("consul.port","8500","")
	)
   flag.Parse()

	//创建日志组件
	//log.SetFlags(log.Ldate|log.Lmicroseconds)
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		//打印时间 装饰器
		logger = log.With(logger,"ts",log.DefaultTimestamp)
		//打印调用对象  装饰器
		logger = log.With(logger,"caller",log.DefaultCaller)
	}

	// 创建consul api客户端
	consulConfig := api.DefaultConfig()
	consulConfig.Address = "http://" + *consulHost+":"+*consulPort
    consulClient,err := api.NewClient(consulConfig)
    if err !=nil {
    	logger.Log("err",err)
    	os.Exit(1)
	}

	//创建反向代理
	proxy := NewReverseProxy(consulClient,logger)

	//信号处理
	errc := make(chan error)
	go func(){
		c := make(chan os.Signal)
		signal.Notify(c,syscall.SIGINT,syscall.SIGTERM)
		errc <-  fmt.Errorf("%s",<-c)
	}()


	//开始监听
	go func() {
		logger.Log("transport", "HTTP", "addr", "9090")
	    errc <-	http.ListenAndServe(":8880",proxy)
	}()

	// 开始运行，等待结束
	logger.Log("exit", <-errc)
}

func NewReverseProxy(client *api.Client,logger log.Logger) *httputil.ReverseProxy {
	//创建Direct
	director := func(req *http.Request)  {
		//查询原始请求路径，如：/arithmetic/calculate/10/5
		reqPath := req.URL.Path
		if reqPath==""{
			return
		}

		//按照分隔符'/'对路径进行分解，获取服务名称serviceName
		pathArray := strings.Split(reqPath,"/")
		serviceName := pathArray[1]

		//调用consul api查询serviceName的服务实例列表
		result,_,err := client.Catalog().Service(serviceName,"",nil)
		if err!=nil {
			logger.Log("ReverseProxy failed", "query service instace error", err.Error())
			return
		}

		if len(result) == 0 {
			logger.Log("ReverseProxy failed", "no such service instance", serviceName)
			return
		}

		//重新组织请求路径 去掉服务部分
		destPath := strings.Join(pathArray[2:],"/")

		//随机选择一个服务实例
		tgt := result[rand.Int()%len(result)]
		logger.Log("service id", tgt.ServiceID)

		req.URL.Scheme = "http"
		req.URL.Host = fmt.Sprintf("%s:%d",tgt.ServiceAddress,tgt.ServicePort)
		req.URL.Path = "/"+destPath
	}
	//Director func(*http.Request) 构造新的 *http.Request
	return &httputil.ReverseProxy{Director:director}
}