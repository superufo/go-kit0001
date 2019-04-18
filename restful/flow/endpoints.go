package flow

import (
	"context"
	"errors"
	"strings"
	"github.com/go-kit/kit/endpoint"
)

/*
在gokit中Endpoint是可以包装到http.Handler中的特殊方法，gokit采用装饰着模式，
把Service应该执行的逻辑封装到Endpoint方法中执行。
Endpoint的作用是：调用Service中相应的方法处理请求对象（ArithmeticRequest），
返回响应对象（ArithmeticResponse）
 */
var  (
	ErrInvalidRequestType = errors.New("RequestType has only four type: Add,Subtract,Multiply,Divide")
)

type ArithmeticRequest struct{
	RequestType string `json:"request_type"`
	A           int `json:"a"`
	B 			int  `json:"b"`
}

// ArithmeticResponse define response struct
type ArithmeticResponse struct {
	Result int   `json:"result"`
	Error  error `json:"error"`
}


func MakeArithmeticEndpoint(svc Service) endpoint.Endpoint{
	return func (ctx context.Context ,request interface{})(response interface{},err error){
			req := request.(ArithmeticRequest)

			var (
				res,a,b int
				calError error
			)

            a = req.A
            b = req.B

			if strings.EqualFold(req.RequestType,"Add"){
				res = svc.Add(a, b)
			}else if strings.EqualFold(req.RequestType,"Substract"){
				res = svc.Subtract(a, b)
			}else if strings.EqualFold(req.RequestType,"Multiply"){
				res = svc.Multiply(a, b)
			}else if strings.EqualFold(req.RequestType,"Divide"){
				res, calError = svc.Divide(a, b)
			}else{
				return nil, ErrInvalidRequestType
			}


            return ArithmeticResponse{Result: res, Error: calError}, nil
	}
}

