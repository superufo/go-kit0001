package main

import  (
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"net/http"
	"os"
	"os/signal"
	flow2 "ratelimit/flow"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"

	"ratelimit/flow"
)

func main(){
	ctx := context.Background()
	errChan := make(chan error)

	var svc  flow.Service
	svc = flow.ArithmeticService{}

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	//包裹日志 进行装饰
	svc = flow.LogginMiddleware(logger)(svc)
	//传入请求 ，传出响应（）， 调用service 处理业务逻辑
	endpoint := flow.MakeArithmeticEndpoint(svc)

	//包裹 限流器 判断流量是否溢出
	// add ratelimit,refill every second,set capacity 3
	//ratebucket := ratelimit.NewBucket(time.Second*1, 3)
	//endpoint = NewTokenBucketLimitterWithJuju(ratebucket)(endpoint)
	// 一秒3次的流量  怎么传进来  就怎么传出去  中间做了限流的异常判断
	ratelimitter := rate.NewLimiter(rate.Every(time.Second*1),3)
    endpoint = flow2.NewTokenBucketLimitterWithBuildIn(ratelimitter)(endpoint)

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
