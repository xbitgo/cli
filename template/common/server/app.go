package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"

	"github.com/xbitgo/core/log"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/sync/errgroup"

	"{{ProjectName}}/common/cfg"
)

type Svc interface {
	Start() error
	Stop()
	StartAddr() string
	Type() string
}

type App struct {
	cfg        *cfg.Server
	svcList    []Svc
	closeFunc  func()
	cancelFunc context.CancelFunc
}

func NewApp(cfg *cfg.Server) *App {
	return &App{
		cfg:     cfg,
		svcList: make([]Svc, 0),
	}
}

func (a *App) Start(etcd *etcdv3.Client) error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	group, ctx := errgroup.WithContext(ctx)
	a.cancelFunc = cancelFunc

	for _, svc := range a.svcList {
		_svc := svc
		group.Go(func() error {
			if err := _svc.Start(); err != nil {
				log.Panicf("starting %s server, err: %s", _svc.Type(), err)
				panic(err)
			}
			return nil
		})
		group.Go(func() error {
			addrName := fmt.Sprintf("%s/%s", a.cfg.AddrName, _svc.Type())
			if err := register(etcd, ctx, addrName, _svc.StartAddr()); err != nil {
				log.Panicf("register %s server, err: %s", _svc.Type(), err)
				panic(err)
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

func ExposedAddr(addr net.Addr) net.Addr {
	tcpAddr, ok := addr.(*net.TCPAddr)
	if !ok {
		return addr
	}
	if !tcpAddr.IP.IsUnspecified() {
		return addr
	}
	_, err := os.Hostname()
	if err != nil {
		log.Errorf("ServiceRegistry convertTcpAddr Hostname fail, err:%s", err)
		return addr
	}
	ips := make([]net.IP, 0)
	interfaceAddrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Errorf("ServiceRegistry convertTcpAddr InterfaceAddrs fail, err:%s", err)
		return addr
	}
	ips = make([]net.IP, 0)
	for _, addr := range interfaceAddrs {
		if ipAddr, ok := addr.(*net.IPNet); ok {
			ips = append(ips, ipAddr.IP)
		} else {
			log.Debugf("bot ip addr addr:[%s], type:[%s]", addr, reflect.TypeOf(addr))
		}
	}
	found := false
	for _, ip := range ips {
		if ip.IsUnspecified() || ip.IsLoopback() {
			continue
		}
		if !strings.HasPrefix(ip.String(), "10.") &&
			!strings.HasPrefix(ip.String(), "192.168.") &&
			!strings.HasPrefix(ip.String(), "172.") {
			continue
		}
		tcpAddr.IP = ip
		found = true
		break
	}
	if !found {
		log.Errorf("ServiceRegistry convertTcpAddr LookupIP no suitable ip found, host:%s", ips)
		return addr
	}
	return tcpAddr
}
