package rpc_client

import (
	"context"
	"time"

	"github.com/pkg/errors"
	etcd "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	RpcConn = map[string]*grpc.ClientConn{}
	From    = ""
)

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
	opt = append(opt,
		grpc.WithResolvers(etcdResolver),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithBlock(),
	)
	return newGrpcClientConn(serviceName, opt...)
}

func newGrpcClientConn(serviceName string, opt ...grpc.DialOption) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	opt = append(opt, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
		//middleware.GRPCClientAuthContext(),
		//middleware.GrpcClientTrace(),
		),
	)
	conn, err := grpc.DialContext(
		ctx,
		serviceName,
		opt...,
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
