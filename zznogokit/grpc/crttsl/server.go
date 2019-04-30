package main

import (
	"crttsl/pb"
	"crttsl/service"
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"net"
)

var (
	port = ":5000"

	tslDir = "./config"

	ca = tslDir + "/ca.crt"
	server_crt = tslDir + "/server.crt"
	server_key = tslDir + "/server.key"
)

func main(){
	certificate,err := tls.LoadX509KeyPair(server_crt,server_key)
	if err !=nil {
		log.Panicf("could not load server key pair :%s",err)
	}

	certPool := x509.NewCertPool()
	ca, err  := ioutil.ReadFile(ca)
	if err !=nil {
		log.Panicf("could not read ca certificate :%s",err)
	}

	if ok:= certPool.AppendCertsFromPEM(ca);!ok{
		log.Panicf("failure to append client certs")
	}

	creds := credentials.NewTLS(&tls.Config{
		ClientAuth:tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs: certPool,
	})

	server := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterGreeterServer(server,new(service.MyGrpcServer))

	lis ,err := net.Listen("tcp",port)
	if err !=nil {
		log.Panicf("could not list on %s: %s", port, err)
	}

	if err := server.Serve(lis); err != nil {
		log.Panicf("grpc serve error: %s", err)
	}
}