package rpc_client

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/xbitgo/components/tracing"
	etcd "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var (
	RpcConn = map[string]*grpc.ClientConn{}
	From    = ""
)

type RPCService struct {
	Name string
	Addr string
}

func RegisterRPC(fromSrv string, etcdClient *etcd.Client, services ...RPCService) (err error) {
	From = fromSrv
	for _, service := range services {
		addr := service.Addr
		if addr == "" {
			return errors.Errorf("can not found service[%s] addr", service.Name)
		}
		conn, err := getConn(etcdClient, addr)
		if err != nil {
			return errors.Wrap(err, "service:"+service.Name+":"+addr)
		}
		RpcConn[service.Name] = conn
	}
	return nil
}

func getConn(etcdClient *etcd.Client, serviceName string, opt ...grpc.DialOption) (*grpc.ClientConn, error) {
	etcdResolver, err := resolver.NewBuilder(etcdClient)
	if err != nil {
		return nil, err
	}
	target := fmt.Sprintf("%s:///%s", etcdResolver.Scheme(), serviceName)
	opt = append(opt,
		grpc.WithResolvers(etcdResolver),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithBlock(),
	)
	return newGrpcClientConn(target, opt...)
}

func newGrpcClientConn(target string, opt ...grpc.DialOption) (*grpc.ClientConn, error) {
	opt = append(opt, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			clientMetadataContext(),
			tracing.GrpcClientTrace(),
		),
	)
	conn, err := grpc.Dial(
		target,
		opt...,
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func clientMetadataContext() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		incomingContext, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			ctx = metadata.NewOutgoingContext(ctx, metadata.MD{})
		}
		ctx = metadata.NewOutgoingContext(ctx, incomingContext)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
