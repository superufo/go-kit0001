package main

import (
	"io"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	. "rpcf/services"
)

func main() {
	rpc.RegisterName("HelloServicek", new(HelloServicek))
	http.HandleFunc("/jsonrpc", func(w http.ResponseWriter, r *http.Request) {
		var conn io.ReadWriteCloser = struct {
			io.Writer
			io.ReadCloser
		}{
			ReadCloser: r.Body,
			Writer: w,
		}
		rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
	})
	http.ListenAndServe(":1234", nil)
}
