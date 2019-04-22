package flow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/gorilla/mux"
	"github.com/go-kit/kit/tracing/zipkin"
	gozipkin "github.com/openzipkin/zipkin-go"
)

/******************************************************
Transport层用于接收用户网络请求并将其转为Endpoint可以处理的对象，
然后交由Endpoint执行，最后将处理结果转为响应对象向用户响应。
为了完成这项工作，Transport需要具备两个工具方法：

解码器：把用户的请求内容转换为请求对象（ArithmeticRequest）；
编码器：把处理结果转换为响应对象（ArithmeticResponse）；
********************************
gorilla/mux是一个强大的路由，小巧但是稳定高效，
不仅可以支持正则路由还可以按照Method，header，host等信息匹配
，可以从我们设定的路由表达式中提取出参数方便上层应用，而且完全兼容http.ServerMux
***************************************************************
*/

var (
	ErrorBadRequest = errors.New("invalid request parameter")
)


// MakeHttpHandler make http handler use mux
func MakeHttpHandler(ctx context.Context, endpoint SetsEndpoints, zipkinTracer *gozipkin.Tracer,logger log.Logger) http.Handler {
	r := mux.NewRouter()

	zipkinServer := zipkin.HTTPServerTrace(zipkinTracer, zipkin.Name("http-transport"))

	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(kithttp.DefaultErrorEncoder),
		zipkinServer,
	}

	fmt.Println(10000)
	r.Methods("POST").Path("/calculate/{type}/{a}/{b}").Handler(kithttp.NewServer(
		endpoint.ArithmeticEndpoint,
		decodeArithmeticRequest,
		encodeArithmeticResponse,
		options...,
	))

	fmt.Println("request:",r)

	//api 监控
	r.Path("/metrics").Handler(promhttp.Handler())

	// 健康检查
	r.Methods("GET").Path("/health").Handler(kithttp.NewServer(
		endpoint.HealthCheckEndpoint,
		decodeHealthCheckRequest,
		encodeArithmeticResponse,
		options...,
	))

	return r
}

// decodeArithmeticRequest decode requestrams to struct
func decodeArithmeticRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	requestType, ok := vars["type"]
	//DefaultErrorEncoder 对应编码
	if !ok {
		return nil, ErrorBadRequest
	}

	pa, ok := vars["a"]
	if !ok {
		return nil, ErrorBadRequest
	}

	pb, ok := vars["b"]
	if !ok {
		return nil, ErrorBadRequest
	}

	a, _ := strconv.Atoi(pa)
	b, _ := strconv.Atoi(pb)

	return ArithmeticRequest{
		RequestType: requestType,
		A:           a,
		B:           b,
	}, nil
}

// encodeArithmeticResponse encode response to return
func encodeArithmeticResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// decodeHealthCheckRequest decode request
func decodeHealthCheckRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return HealthRequest{}, nil
}
