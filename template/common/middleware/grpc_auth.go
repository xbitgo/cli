package middleware

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func GRPCServerAuthContext() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.MD{}
		}
		ctx = metadata.NewIncomingContext(ctx, md)
		return handler(ctx, req)
	}
}

func GRPCClientAuthContext() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		incomingContext, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			ctx = metadata.NewOutgoingContext(ctx, metadata.MD{})
		}
		ctx = metadata.NewOutgoingContext(ctx, incomingContext)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
