package servers

import "net/rpc"

//名字
const HelloServiceName = "servers/test/HelloService"

//接口  要暴露的方法
type HelloServiceInterface  interface {
	Hello(request string,reply *string) error
}

//type HelloService struct{}

//服务  struct
func RegistertHelloService(svc HelloServiceInterface) error{
	return rpc.RegisterName(HelloServiceName,svc)
}

