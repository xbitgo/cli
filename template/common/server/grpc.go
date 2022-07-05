package server

import (
	"fmt"
	"net"
	"runtime/debug"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xbitgo/core/log"
)

var (
	grpcRecoveryOpts = []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			log.Errorf("grpc_recovery: %v \n", p)
			debug.PrintStack()
			return status.Errorf(codes.Unknown, "panic triggered: %v", p)
		}),
	}

	grpcUnaryServerOpts = []grpc.UnaryServerInterceptor{
		grpc_recovery.UnaryServerInterceptor(grpcRecoveryOpts...),
		grpc_prometheus.UnaryServerInterceptor,
	}
)

type grpcSvc struct {
	Listener net.Listener
	svc      *grpc.Server
}

func (s *grpcSvc) Type() string {
	return "grpc"
}

func (s *grpcSvc) StartAddr() string {
	return ExposedAddr(s.Listener.Addr()).String()
}

func (s *grpcSvc) Start() error {
	fmt.Printf(" -- starting grpc server: [%s] ... \n", s.Listener.Addr())
	return s.svc.Serve(s.Listener)
}

func (s *grpcSvc) Stop() {
	s.svc.GracefulStop()
}

func (a *App) InitGRPC(process func(grpcService grpc.ServiceRegistrar), opts ...grpc.UnaryServerInterceptor) {
	addr := a.cfg.GRPCAddr
	if addr == "" {
		return
	}
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(fmt.Errorf("-- tcp listening on addr: %s, err: %v \n ", addr, err))
		return
	}
	// 配置自定义拦截器
	grpcUnaryServerOpts = append(grpcUnaryServerOpts, opts...)
	grpcOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpcUnaryServerOpts...,
		)),
	}
	grpcServer := grpc.NewServer(grpcOptions...)
	process(grpcServer)
	a.svcList = append(a.svcList, &grpcSvc{
		Listener: lis,
		svc:      grpcServer,
	})
	return
}
