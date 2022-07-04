package cfg

import (
	etcdv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type Etcd struct {
	Endpoints []string `json:"endpoints" yaml:"endpoints"`
	Timeout   int      `json:"timeout" yaml:"timeout"`
}

// CreateInstance DI生成器默认调用
func (e *Etcd) CreateInstance() (*etcdv3.Client, error) {
	timeout := 5 * time.Second
	if e.Timeout != 0 {
		timeout = time.Duration(e.Timeout) * time.Second
	}
	cli, err := etcdv3.New(etcdv3.Config{
		Endpoints:   e.Endpoints,
		DialTimeout: timeout,
	})
	if err != nil {
		return nil, err
	}
	return cli, nil
}
