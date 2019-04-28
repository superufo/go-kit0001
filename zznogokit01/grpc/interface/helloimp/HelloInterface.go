package helloimp

import (
	"context"
	"google.golang.org/grpc"
	"strings"
)

type HelloServiceServer interface {
	Hello(context.Context, *String) (*String, error)
}
type HelloServiceClient interface {
	Hello(context.Context, *String, ...grpc.CallOption) (*String
, error)
}
