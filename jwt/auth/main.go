package main

//	"gopkg.in/alecthomas/kingpin.v2"
import  (
	flow "auth/flow"
	"context"
	"flag"
	"fmt"
	"golang.org/x/time/rate"

	"net/http"
	"os"
	"os/signal"

	"syscall"
	"time"

	"github.com/go-kit/kit/log"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

func main(){
	var (
		consulHost  = flag.String("consul.host", "", "consul ip address")
		consulPort  = flag.String("consul.port", "", "consul port")
		serviceHost = flag.String("service.host", "", "service ip address")
		servicePort = flag.String("service.port", "", "service port")
		zipkinURL   = flag.String("zipkin.url", "http://192.168.192.146:9411/api/v2/spans", "Zipkin server url")
	)
	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	fieldKeys := []string{"method"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace:"mike",
		Subsystem:"arithmetic_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	},fieldKeys)

	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "mike",
		Subsystem: "arithemetic_service",
		Name:      "request_latency",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	var zipkinTracer *zipkin.Tracer
	{
		var (
			err           error
			hostPort      = *serviceHost + ":" + *servicePort
			serviceName   = "arithmetic-service"
			useNoopTracer = (*zipkinURL == "")
			reporter      =  zipkinhttp.NewReporter(*zipkinURL)
		)
		defer reporter.Close()
		zEP, _ := zipkin.NewEndpoint(serviceName, hostPort)
		zipkinTracer, err = zipkin.NewTracer(
			reporter, zipkin.WithLocalEndpoint(zEP), zipkin.WithNoopTracer(useNoopTracer),
		)
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		if !useNoopTracer {
			logger.Log("tracer", "Zipkin", "type", "Native", "URL", *zipkinURL)
		}
	}



	var svc  flow.Service
	svc = flow.ArithmeticService{}

	//包裹日志 进行装饰
	svc = flow.LogginMiddleware(logger)(svc)

	//Api监控 进行装饰
	svc = flow.Metrics(requestCount, requestLatency)(svc)


	//包裹 限流器 判断流量是否溢出
	// add ratelimit,refill every second,set capacity 3
	//ratebucket := ratelimit.NewBucket(time.Second*1, 3)
	//endpoint = NewTokenBucketLimitterWithJuju(ratebucket)(endpoint)
	// 一秒3次的流量  怎么传进来  就怎么传出去  中间做了限流的异常判断
	ratelimitter := rate.NewLimiter(rate.Every(time.Second*1),100)

	//传入请求 ，传出响应（）， 调用service 处理业务逻辑
	endpoint := flow.MakeArithmeticEndpoint(svc)
    endpoint = flow.NewTokenBucketLimitterWithBuildIn(ratelimitter)(endpoint)
	endpoint = kitzipkin.TraceEndpoint(zipkinTracer, "calculate-endpoint")(endpoint)

    /*************** 健康检查 ***********************/
    //创建健康检查的Endpoint，未增加限流
    healthEndpoint := flow.MakeHealthCheckEndpoint(svc)
	healthEndpoint = flow.NewTokenBucketLimitterWithBuildIn(ratelimitter)(healthEndpoint)
	healthEndpoint = kitzipkin.TraceEndpoint(zipkinTracer, "health-endpoint")(healthEndpoint)

	//身份认证 Endpoint
	authEndpoint := flow.MakeAuthEndpoint(svc)
	authEndpoint  = flow.NewTokenBucketLimitterWithBuildIn(ratelimitter)(authEndpoint)
	authEndpoint  = kitzipkin.TraceEndpoint(zipkinTracer, "login-endpoint")(authEndpoint)


	//把算术运算Endpoint和健康检查Endpoint封装至ArithmeticEndpoints
	endpts := flow.SetsEndpoints{
		ArithmeticEndpoint:  endpoint,
		HealthCheckEndpoint: healthEndpoint,
		AuthEndpoint : authEndpoint,
	}

	//创建http.Handler 处理传输协议 （http,tcp grpc等） 处理返回的格式
	r := flow.MakeHttpHandler(ctx, endpts,zipkinTracer, logger)

	//创建注册对象
	registar := flow.Register(*consulHost, *consulPort, *serviceHost, *servicePort, logger)

	go func() {
		fmt.Println("Http Server start at port:" + *servicePort)
		//启动前执行注册
		registar.Register()
		handler := r
		errChan <- http.ListenAndServe(":"+*servicePort, handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errChan

	//服务退出取消注册
	registar.Deregister()
	fmt.Println(error)
}