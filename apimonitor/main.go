package main

import  (
	"context"
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
	"apimonitor/flow"
)

func main(){
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
		Namespace:"raysonxin",
		Subsystem:"arithmetic_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	},fieldKeys)

	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "raysonxin",
		Subsystem: "arithemetic_service",
		Name:      "request_latency",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	var svc  flow.Service
	svc = flow.ArithmeticService{}

	//包裹日志 进行装饰
	svc = flow.LogginMiddleware(logger)(svc)

	//Api监控 进行装饰
	svc = flow.Metrics(requestCount, requestLatency)(svc)

	//传入请求 ，传出响应（）， 调用service 处理业务逻辑
	endpoint := flow.MakeArithmeticEndpoint(svc)


	//包裹 限流器 判断流量是否溢出
	// add ratelimit,refill every second,set capacity 3
	//ratebucket := ratelimit.NewBucket(time.Second*1, 3)
	//endpoint = NewTokenBucketLimitterWithJuju(ratebucket)(endpoint)
	// 一秒3次的流量  怎么传进来  就怎么传出去  中间做了限流的异常判断
	ratelimitter := rate.NewLimiter(rate.Every(time.Second*1),100)
    endpoint = flow.NewTokenBucketLimitterWithBuildIn(ratelimitter)(endpoint)

	//处理传输协议 （http,tcp grpc等） 处理返回的格式
	r := flow.MakeHttpHandler(ctx, endpoint, logger)

	go func() {
		fmt.Println("Http Server start at port:9550")
		handle := r
		errChan <- http.ListenAndServe("127.0.0.1:9550",handle)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	fmt.Println(<-errChan)
}
