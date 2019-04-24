package servers

import "net/rpc"

type HelloServiceClient struct {
	*rpc.Client
}

//转换成接口
var _ HelloServiceInterface = (*HelloServiceClient)(nil)

func DialHelloService (network,address string)(*HelloServiceClient,error){
	c,err := rpc.Dial(network,address)
	if err!=nil {
		return nil,err
	}
	return &HelloServiceClient{Client:c},nil
}

func (p *HelloServiceClient)Hello(request string,reply *string) error {
	//*reply =  "hello world, request:" +  request
	//return nil
	return p.Client.Call(HelloServiceName+".Hello",request,reply)
}


