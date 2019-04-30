package main

import (
	"fmt"
	"net/http"
	"strings"
	"webtsl/pb"
	"webtsl/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

var (
	port   = ":5000"
	tslDir = "./config"
	//tlsServerName = "server.grpc.io"

	serverCrt = tslDir + "/server.crt"
	serverKey = tslDir + "/server.key"
)

func main() {
	creds, err := credentials.NewServerTLSFromFile(serverCrt, serverKey)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterGreeterServer(grpcServer, new(service.MyGrpcServer))

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(w, "hello")
	})

	http.ListenAndServeTLS(port, serverCrt, serverKey, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO(tamird): point to merged gRPC code rather than a PR.
		// This is a partial recreation of gRPC's internal checks
		// https://github.com/grpc/grpc-go/pull/514/files#diff-95e9a25b738459a2d3030e1e6fa2a718R61
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			mux.ServeHTTP(w, r)
		}
	}))

}
