package main

import (
	"errors"
	"fmt"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/log"
	"github.com/hashicorp/consul/api"
	"github.com/openzipkin/zipkin-go"
	zipkinhttpsvr "github.com/openzipkin/zipkin-go/middleware/http"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
)

//HystrixRouter hystrix 路由
type HystrixRouter struct {
	svcMap   *sync.Map   //服务实例 存储已经通过HYstrix 监控服务列表
	logger   log.Logger  //日志工具
    fallbackMsg  string //回调信息
    consulClient *api.Client   //consule 客户端对象
    tracer   *zipkin.Tracer
}

func Routes(client *api.Client,zikkinTracer *zipkin.Tracer, fbMsg string, logger log.Logger) http.Handler {
	return HystrixRouter{
		svcMap:       &sync.Map{},
		logger:       logger,
		fallbackMsg:  fbMsg,
		consulClient: client,
		tracer:       zikkinTracer,
	}
}

func (route HystrixRouter) ServerHTTP(h http.ResponseWriter,r *http.Request){
	//查询原始请求路径，如：/arithmetic/calculate/10/5
	reqPath := r.URL.Path
	if reqPath==""{
		return
	}

	//按照分隔符'/'对路径进行分解，获取服务名称serviceName
	pathArray := strings.Split(reqPath,"/")
    serviceName := pathArray[1]

    //检查是否已经加入监控
	if _,ok := route.svcMap.Load(serviceName);!ok {
		//吧serverName 作为命令对象，设置参数
		hystrix.ConfigureCommand(serviceName,hystrix.CommandConfig{Timeout:1000})
		route.svcMap.Store(serviceName,serviceName)
	}

	//hystrix.Do("", func()error{},func(err error)error{})
	//func()error{} 正常业务逻辑，一般是访问其他资源
	//func(err error)error{}  失败处理逻辑，访问其他资源失败时，或者处于熔断开启状态时，会调用这段逻辑
	//可以简单构造一个response返回，也可以有一定的策略，比如访问备份资源
	//也可以直接返回 err，这样不用和远端失败的资源通信，防止雪崩
	err:= hystrix.Do(serviceName, func() error {
	//调用consul api查询serviceNam
	result,_,err := route.consulClient.Catalog().Service(serviceName,"",nil)
    if err!=nil {
		router.logger.Log("ReverseProxy failed", "query service instace error", err.Error())
		return
	}

	if len(result)==0 {
		router.logger.Log("ReverseProxy failed", "no such service instance", serviceName)
		return errors.New("no such service instance")
	}

	director := func(req *http.Request){
		//重新组织请求路径，去掉服务名称部分
		destPath := strings.Join(pathArray[2:], "/")



	}






	})


}









