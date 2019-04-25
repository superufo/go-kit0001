package service

import (
	"fmt"
	"log"
	"net"
)

type Service interface {
	Hello(request string,reply *string) error
	Login(request string,reply *string) error
}

type HelloService struct{
	conn net.Conn
	isLogin bool
}

func (p *HelloService) Hello (request string ,reply *string) error {
	if !p.isLogin {
		return fmt.Errorf("please login")
	}

	*reply = "hello" + request + ",from " + p.conn.RemoteAddr()
	return nil
}

func (p *HelloService) Login(request string,reply *string) error{
	if request !="user:password" {
		return  fmt.Errorf("auth failed")
	}

	log.Println("login ok")
	p.isLogin = true
	return nil
}