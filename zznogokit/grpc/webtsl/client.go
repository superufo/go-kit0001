package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"webtsl/pb"
)

var (
	port          = ":5000"
	tslDir        = "./config"
	tlsServerName = "server.grpc.io"

	serverCrt = tslDir + "/server.crt"
)

func main() {
	creds, err := credentials.NewClientTLSFromFile(serverCrt, "server.grpc.io")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial("localhost"+port, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "gopher"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("doClientWork: %s", r.Message)
}
