package flow

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	kithttp "github.com/go-kit/kit/transport/http"

	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

/* *********************************************************************************
*  kit/sd/Endpointer:
*  Endpointer listens to a service discovery system and yields a set of
*  identical endpoints on demand. An error indicates a problem with connectivity
*  to the service discovery system, or within the system itself; an Endpointer
*  may yield no endpoints without error.
*  type Endpointer interface {
*    Endpoints() ([]endpoint.Endpoint, error)
*  }
*
*  NewEndpointer creates an Endpointer that subscribes to updates from Instancer src
*  and uses factory f to create Endpoints. If src notifies of an error, the Endpointer
*  keeps returning previously created Endpoints assuming they are still good, unless
*  this behavior is disabled via InvalidateOnError option.
*  func NewEndpointer(src Instancer, f Factory, logger log.Logger, options ...EndpointerOption) *DefaultEndpointer
*  kit/sd/Endpointer提供了一套服务发现机制，其定义和创建接口如下所
*  Endpointer通过监听服务发现系统的事件信息，并且通过factory按需创建服务终结点（Endpoint）。
*  所以，我们需要通过Endpointer来实现服务发现功能。在微服务模式下，同一个服务可能存在多个实例，
*  所以需要通过负载均衡机制完成实例选择，这里使用go-kit工具集中的kit/sd/lb组件（该组件实现RoundRibbon，并具备Retry功能）
*  *********************************************************************************/

func arithmeticFactory(_ context.Context,method,path string) sd.Factory{
	return func(instance string)(endpoint endpoint.Endpoint,closer io.Closer,err error) {
		if !strings.HasPrefix(instance, "http" ) {
		  instance = "http://" + instance
		}

		tgt ,err := url.Parse(instance)
		if err!=nil {
			return nil,nil,err
		}

		tgt.Path = path

		var (
			enc kithttp.EncodeRequestFunc
			dec kithttp.DecodeResponseFunc
		)

		enc , dec = encodeArithmeticRequest, decodeArithmeticReponse

		return kithttp.NewClient(method,tgt,enc,dec).Endpoint(),nil,nil
	}
}

func encodeArithmeticRequest(_ context.Context, req *http.Request, request interface{}) error {
	arithReq := request.(ArithmeticRequest)
	p := "/" + arithReq.RequestType + "/" + strconv.Itoa(arithReq.A) + "/" + strconv.Itoa(arithReq.B)
	req.URL.Path += p
	return nil
}

func decodeArithmeticReponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response ArithmeticResponse
	var s map[string]interface{}

	if respCode := resp.StatusCode; respCode >= 400 {
		if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
			return nil, err
		}
		return nil, errors.New(s["error"].(string) + "\n")
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}