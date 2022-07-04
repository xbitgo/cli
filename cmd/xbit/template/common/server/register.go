package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/xbitgo/core/tools/tool_str"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

const (
	prefix = "service"
)

func Register(etcd *clientv3.Client, ctx context.Context, serviceName, addr string) error {
	log.Println("Try register to etcd ...")
	// 创建一个租约
	lease := clientv3.NewLease(etcd)
	cancelCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	leaseResp, err := lease.Grant(cancelCtx, 3)
	if err != nil {
		return err
	}

	leaseChannel, err := lease.KeepAlive(ctx, leaseResp.ID) // 长链接, 不用设置超时时间
	if err != nil {
		return err
	}

	em, err := endpoints.NewManager(etcd, prefix)
	if err != nil {
		return err
	}

	cancelCtx, cancel = context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	if err := em.AddEndpoint(cancelCtx, fmt.Sprintf("%s/%s/%s", prefix, serviceName, tool_str.UUID()), endpoints.Endpoint{
		Addr: addr,
	}, clientv3.WithLease(leaseResp.ID)); err != nil {
		return err
	}
	log.Println("Register etcd success")

	del := func() {
		log.Println("Register close")
		cancelCtx, cancel = context.WithTimeout(ctx, time.Second*3)
		defer cancel()
		em.DeleteEndpoint(cancelCtx, serviceName)
		lease.Close()
	}
	// 保持注册状态(连接断开重连)
	keepRegister(etcd, ctx, leaseChannel, del, serviceName, addr)

	return nil
}

func keepRegister(etcd *clientv3.Client, ctx context.Context, leaseChannel <-chan *clientv3.LeaseKeepAliveResponse, cleanFunc func(), serviceName, addr string) {
	go func() {
		failedCount := 0
		for {
			select {
			case resp := <-leaseChannel:
				if resp != nil {
					//log.Println("keep alive success.")
				} else {
					log.Println("keep alive failed.")
					failedCount++
					for failedCount > 3 {
						cleanFunc()
						if err := Register(etcd, ctx, serviceName, addr); err != nil {
							time.Sleep(time.Second)
							continue
						}
						return
					}
					continue
				}
			case <-ctx.Done():
				cleanFunc()
				etcd.Close()
				return
			}
		}
	}()
}
