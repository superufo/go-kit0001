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
	"github.com/openzipkin/zipkin-go"
	zipkinhttpsvr "github.com/openzipkin/zipkin-go/middleware/http"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

func main(){
	// 创建环境变量
	var (
		consulHost = flag.String("consul.host","47.112.111.171","")
		consulPort = flag.String("consul.port","8500","")
		zipkinURL = flag.String("zipkin.url","http://47.112.111.171:9411/zipkin/api/v2/spans","zipnkin server url")
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

	var zipkinTracer  *zipkin.Tracer
	{
		var (
			err error
			hostPort = "127.0.0.1:8880"
			serviceName = "gateway-service"
			useNoopTracer = (*zipkinURL == "")
			reporter      = zipkinhttp.NewReporter(*zipkinURL)
		)

		defer reporter.Close()

		zep,_ := zipkin.NewEndpoint(serviceName,hostPort)
		zipkinTracer, err = zipkin.NewTracer(
			reporter,zipkin.WithLocalEndpoint(zep),zipkin.WithNoopTracer(useNoopTracer),
			)
		if err != nil {
			logger.Log("err",err)
		}

		if !useNoopTracer{
			logger.Log("tracer", "Zipkin", "type", "Native", "URL", *zipkinURL)
		}

	 }

	// 创建consul api客户端
	consulConfig := api.DefaultConfig()
	consulConfig.Address = "http://" + *consulHost + ":" + *consulPort
	consulClient, err := api.NewClient(consulConfig)
    if err !=nil {
    	logger.Log("err",err)
    	os.Exit(1)
	}

	//创建反向代理   zipkinTracer 传递
	proxy := NewReverseProxy(consulClient,zipkinTracer,logger)

	tags := map[string]string{
		"component":"gateway_server",
	}

	//NewServerMiddleware(t *zipkin.Tracer, options ...ServerOption) func(http.Handler) http.Handler
	//装饰器包裹
	handle := zipkinhttpsvr.NewServerMiddleware(
		zipkinTracer,
		zipkinhttpsvr.SpanName("gateway"),
		zipkinhttpsvr.TagResponseSize(true),
		zipkinhttpsvr.ServerTags(tags),
		)(proxy)



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
	    errc <-	http.ListenAndServe(":8880",handle)
	}()

	// 开始运行，等待结束
	logger.Log("exit", <-errc)
}

func NewReverseProxy(client *api.Client, zikkinTracer *zipkin.Tracer,logger log.Logger) *httputil.ReverseProxy {
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

	/**gateway接收请求后，会创建一个span，其中的traceId将作为本次请求的唯一编号，
	gateway必须把这个traceId“告诉”算术运算服务，算术运算服务才能为该请求持续记录追踪信息。
	在ReverseProxy中能够完成这一任务的就是Transport，我们可以使用zipkin-go的middleware/http包
	提供的NewTransport替换系统默认的http.DefaultTransport
	***/
	rounndTrip, _ := zipkinhttpsvr.NewTransport(zikkinTracer,zipkinhttpsvr.TransportTrace(true))

	//Director func(*http.Request) 构造新的 *http.Request
	return &httputil.ReverseProxy{Director:director,Transport:rounndTrip}
}