package main

import  (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"

	"restful/flow"
)

func main(){
	ctx := context.Background()
	errChan := make(chan error)

	var svc  flow.Service
	svc = flow.ArithmeticService{}

	endpoint := flow.MakeArithmeticEndpoint(svc)
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
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
