package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/xbitgo/core/log"
)

type Svc interface {
	Start() error
	Type() string
}

type App struct {
	grpcSvc *grpcSvc
	httpSvc *httpSvc
	//taskSvc    *taskSvc
	closeFunc  func()
	cancelFunc context.CancelFunc
}

func NewApp() *App {
	return &App{}
}

func (a *App) Start() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	group, ctx := errgroup.WithContext(ctx)
	a.cancelFunc = cancelFunc
	if a.grpcSvc != nil {
		group.Go(func() error {
			if err := a.grpcSvc.Start(); err != nil {
				return fmt.Errorf("starting grpc server, err: %s", err)
			}
			return nil
		})
	}
	if a.httpSvc != nil {
		group.Go(func() error {
			if err := a.httpSvc.Start(); err != nil {
				return fmt.Errorf("starting http server, err: %s", err)
			}
			return nil
		})
	}
	go a.signalExit()
	return group.Wait()
}

func (a *App) OnClose(fun func()) {
	a.closeFunc = fun
}

func (a *App) signalExit() {
	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	for {
		s := <-c
		log.Infof("service get a signal: %v", s)
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT, syscall.SIGHUP:
			a.closeFunc()
			a.cancelFunc()
			log.Info("service closed")
			os.Exit(0)
			return
		default:
			return
		}
	}
}

//func (a *App) StartEtcd(etcdAddr string, etcd *etcdv3.Client) error {
//	if etcd == nil {
//		return errors.New("etcd client err")
//	}
//	etcdGRPCAddr := fmt.Sprintf("%s_grpc", etcdAddr)
//	etcdHTTPAddr := fmt.Sprintf("%s_http", etcdAddr)
//
//	if a.grpcSvc != nil {
//		err := a.grpcSvc.Start()
//		if err != nil {
//			return err
//		}
//	}
//
//	////
//	//if server.GRPCAddr != "" {
//	//	lis, err := net.Listen("tcp", server.GRPCAddr)
//	//	if err != nil {
//	//		log.Errorf("tcp listening on addr: %s, err: %v", server.GRPCAddr, err)
//	//	}
//	//	grpcServer := grpc.NewServer(grpcOptions...)
//	//	defer grpcServer.GracefulStop()
//	//	pb.RegisterLocationServiceServer(grpcServer, loc)
//	//}
//	//if server.HTTPAddr != "" {
//	//
//	//}
//
//	return nil
//}
