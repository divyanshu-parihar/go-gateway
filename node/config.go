package node

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func NewETCDConfig(endpoints []string, timeout time.Duration) (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: timeout,
	})
	return cli, err
}

func NewRedis(addr, password string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{Addr: addr, Password: password})

	return rdb
}

func NodeResourceSetup(ctx context.Context, etcdEndpoints []string, redisCredential string) (*clientv3.Client, *redis.Client, error) {
	cli, err := NewETCDConfig(etcdEndpoints, 3)
	if err != nil {
		slog.Error("ETCD", "status", "Failed to connect to ETCD", "error", err)
		return nil, nil, err
	}

	dns := strings.Split(redisCredential, "+")
	if len(dns) == 1 {
		dns = append(dns, "")
	}
	if dns == nil || (len(dns) != 2) {
		return nil, nil, errors.New("Invalid Redis credentials format. Expected 'host:password'")
	}
	host := dns[0]
	password := dns[1]

	if host == "" {
		return nil, nil, errors.New("Invalid Redis credentials format. Expected 'host:password'")
	}
	rdb := NewRedis(host, password)

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, nil, err
	}
	return cli, rdb, nil
}
