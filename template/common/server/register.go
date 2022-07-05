package server

import (
	"context"
	"fmt"

	"github.com/xbitgo/core/log"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func toAddrVal(addr string) string {
	//return addr
	return fmt.Sprintf(`{"Addr":"%s"}`, addr)
}

func register(etcd *clientv3.Client, ctx context.Context, serviceName, addr string) error {
	leaseResp, err := etcd.Grant(ctx, 3)
	if err != nil {
		return err
	}
	addrName := fmt.Sprintf("%s/%v", serviceName, leaseResp.ID)
	_, err = etcd.Put(ctx, addrName, toAddrVal(addr), clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}
	// 保持注册状态(连接断开重连)
	return keepAlive(etcd, serviceName, addr, leaseResp)
}

func keepAlive(etcd *clientv3.Client, serviceName, addr string, leaseResp *clientv3.LeaseGrantResponse) error {
	leaseChannel, err := etcd.KeepAlive(etcd.Ctx(), leaseResp.ID)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case _, ok := <-leaseChannel:
				if !ok {
					_, err = etcd.Revoke(etcd.Ctx(), leaseResp.ID)
					if err != nil {
						log.Errorf("etcd.Revoke %v; err:%v", leaseResp.ID, err)
					}
					register(etcd, context.Background(), serviceName, addr)
					return
				}
			}
		}
	}()
	return nil
}
