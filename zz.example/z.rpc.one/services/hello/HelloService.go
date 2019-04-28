package hello

/***
Go语言的RPC包的路径为net/rp
中Hello方法必须满足Go语言的RPC规则：方法只能有两个可序列化的参数，其
中第二个参数是指针类型，并且返回一个error类型，同时必须是公开的方法。
 */
type HelloService struct {}

func (p *HelloService) Hello(request string,reply *string) error {
	*reply = "hello:" + request
	return nil
}