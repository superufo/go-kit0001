package flow

import (
	"github.com/go-kit/kit/log"
	"time"
)

type loggingMiddleware struct {
	Service
	logger log.Logger
}

//LoggingMiddleware把日志记录对象嵌入中间件。该方法接受日志对象，
// 返回ServiceMiddleware，而ServiceMiddleware可以传入Service对象，
// 这样就可以对Service增加一层装饰 . 这就是设计模式-装饰者模式
// loggingMiddleware 继承了 service , 传入父类的service,
// 返回继承类的带装饰行为的 service
// 带点 New 的 味道
func LogginMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next Service) Service {
		return loggingMiddleware{next,logger}
	}
}

//包裹父类的方法
func (mw loggingMiddleware)Add (a,b int) (ret int){
	defer func(begin time.Time){
		mw.logger.Log(
			"function", "Add",
						"a", a,
						"b", b,
						"result", ret,
						"took", time.Since(begin),
			)
	}(time.Now())

	ret = mw.Service.Add(a,b)
	return ret
}

func (mw loggingMiddleware) Subtract(a, b int) (ret int) {
	defer func(beign time.Time) {
		mw.logger.Log(
			"function", "Subtract",
			"a", a,
			"b", b,
			"result", ret,
			"took", time.Since(beign),
		)
	}(time.Now())

	ret = mw.Service.Subtract(a, b)
	return ret
}

func (mw loggingMiddleware) Multiply(a, b int) (ret int) {
	defer func(beign time.Time) {
		mw.logger.Log(
			"function", "Multiply",
			"a", a,
			"b", b,
			"result", ret,
			"took", time.Since(beign),
		)
	}(time.Now())

	ret = mw.Service.Multiply(a, b)
	return ret
}

func (mw loggingMiddleware) Divide(a, b int) (ret int, err error) {
	defer func(beign time.Time) {
		mw.logger.Log(
			"function", "Divide",
			"a", a,
			"b", b,
			"result", ret,
			"took", time.Since(beign),
		)
	}(time.Now())

	ret, err = mw.Service.Divide(a, b)
	return
}