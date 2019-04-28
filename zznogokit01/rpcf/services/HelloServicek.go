package services

type HelloServicek struct {}

func (p *HelloServicek) Hello(request string, reply *string) error {
	*reply = "hello:" + request
	return nil
}