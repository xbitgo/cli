package cfg

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type Redis struct {
	Addr     string `json:"addr" yaml:"addr"`
	Password string `json:"password" yaml:"password"`
	PoolSize int    `json:"pool_size" yaml:"pool_size"`
	DB       int    `json:"db" yaml:"db"`
}

func (r *Redis) CreateInstance() (*redis.Client, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:     r.Addr,
		Password: r.Password,
		DB:       r.DB,
		PoolSize: r.PoolSize,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	rs := cli.Ping(ctx)
	return cli, rs.Err()
}
