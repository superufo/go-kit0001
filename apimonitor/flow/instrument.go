package flow

import (
	"context"
	"errors"
	"time"

	"github.com/go-kit/kit/metrics"
	"golang.org/x/time/rate"
	"github.com/go-kit/kit/endpoint"
	"github.com/juju/ratelimit"
)

/********************限流器代码 start**************************************/
var ErrLimitExceed = errors.New("Rate limit excced")

//NewTokenBucketLimitterWithJuju 使用juju/ratelimit创建限流中间件
func NewTokenBucketLimitterWithJuju(bkt *ratelimit.Bucket) endpoint.Middleware{
	return func(next endpoint.Endpoint) endpoint.Endpoint{
		return func(ctx context.Context,request interface{})(response interface{},err error){
			if bkt.TakeAvailable(1) == 0 {
				return nil,ErrLimitExceed
			}

			return next(ctx,request)
		}
	}
}

//NewTokenBucketLimitterWithBuildIn 使用x/time/rate创建限流中间件
func NewTokenBucketLimitterWithBuildIn(bkt *rate.Limiter) endpoint.Middleware{
	return func(next endpoint.Endpoint) endpoint.Endpoint {
  		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !bkt.Allow() {
				return nil, ErrLimitExceed
			}

			return next(ctx,request)
		}
    }
}
/********************限流器代码 end**************************************/

/********************api 监控代码 开始**************************************/
//metricMiddleware 定义监控中间件，嵌入Service 或叫 继承Service接口
// 新增监控指标项：requestCount和requestLatency
type metricMiddleware struct {
	Service
	requestCount  metrics.Counter
	requestLatency metrics.Histogram
}

func Metrics(reqquestCount metrics.Counter,requestLatency metrics.Histogram) ServiceMiddleware{
	return func(next Service)Service{
		return metricMiddleware{
			next,
			reqquestCount,
			requestLatency,
		}
	}
}

/**
可变长参数函数:
原型：func sum(nums ...int)
调用： nums := []int{1, 2, 3, 4}
       sum(nums...)
type Counter interface {
	With(labelValues ...string) Counter
	Add(delta float64)
}
*/
func (mw metricMiddleware) Add(a,b int)(ret int){
	defer func(begin time.Time){
		lvs := []string{"method","Add"}
		//加1的意思
		mw.requestCount.With(lvs...).Add(1)
		// 过去的时间间隔
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	ret = mw.Service.Add(a,b)
	return ret
}

func (mw metricMiddleware) Subtract(a, b int) (ret int) {
	defer func(beign time.Time) {
		lvs := []string{"method", "Subtract"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(beign).Seconds())
	}(time.Now())

	ret = mw.Service.Subtract(a, b)
	return ret
}

func (mw metricMiddleware) Multiply(a, b int) (ret int) {
	defer func(beign time.Time) {
		lvs := []string{"method", "Multiply"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(beign).Seconds())
	}(time.Now())

	ret = mw.Service.Multiply(a, b)
	return ret
}

func (mw metricMiddleware) Divide(a, b int) (ret int, err error) {
	defer func(beign time.Time) {
		lvs := []string{"method", "Divide"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(beign).Seconds())
	}(time.Now())

	ret, err = mw.Service.Divide(a, b)
	return
}