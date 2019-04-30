package main

import (
	"crttsl/pb"
	"crypto/tls"
	"crypto/x509"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"context"
)

var (
	port = ":5000"

	tslDir = "./config"
	tlsServerName = "server.grpc.io"

	ca = tslDir + "/ca.crt"
	client_crt = tslDir + "/client.crt"
	client_key = tslDir + "/client.key"
)

func main(){
	certificate,err := tls.LoadX509KeyPair(client_crt,client_key)
    if err !=nil {
    	log.Panicf("could not load client key pair: %s",err)
	}

	certPool := x509.NewCertPool()
	ca,err := ioutil.ReadFile(ca)
	if err !=nil {
		log.Panicf("could not read ca certificate: %s",err)
	}

	if ok:= certPool.AppendCertsFromPEM(ca);!ok{
		log.Panic("failed to append ca certs")
	}

	creds := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify:false,
		ServerName : tlsServerName,
		Certificates:[]tls.Certificate{certificate},
		RootCAs:certPool,
	})

	conn,err := grpc.Dial("localhost"+port,grpc.WithTransportCredentials(creds))
    if err !=nil {
    	log.Fatal(err)
	}

	defer conn.Close()

	c := pb.NewGreeterClient(conn)
    r,err := c.SayHello(context.Background(),&pb.HelloRequest{Name:"gopher"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("doClientWork: %s", r.Message)
}