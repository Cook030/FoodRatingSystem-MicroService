package registry

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	servicePrefix = "/services/"
	leaseTTL      = 10
)

type EtcdRegistry struct {
	client *clientv3.Client
	lease  clientv3.LeaseID
}

func NewEtcdRegistry(endpoints []string) (*EtcdRegistry, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("连接 etcd 失败: %v", err)
	}
	return &EtcdRegistry{client: cli}, nil
}

func (r *EtcdRegistry) Register(ctx context.Context, serviceName, serviceAddr string) error {
	resp, err := r.client.Grant(ctx, int64(leaseTTL))
	if err != nil {
		return fmt.Errorf("创建租约失败: %v", err)
	}
	r.lease = resp.ID

	key := servicePrefix + serviceName + "/" + serviceAddr
	_, err = r.client.Put(ctx, key, serviceAddr, clientv3.WithLease(r.lease))
	if err != nil {
		return fmt.Errorf("注册服务失败: %v", err)
	}

	keepAliveCh, err := r.client.KeepAlive(ctx, r.lease)
	if err != nil {
		return fmt.Errorf("启动租约续期失败: %v", err)
	}

	go func() {
		for range keepAliveCh {
		}
		log.Printf("[etcd] 服务 %s (%s) 租约续期已停止", serviceName, serviceAddr)
	}()

	log.Printf("[etcd] 服务注册成功: %s -> %s", serviceName, serviceAddr)
	return nil
}

func (r *EtcdRegistry) Deregister(ctx context.Context, serviceName, serviceAddr string) error {
	key := servicePrefix + serviceName + "/" + serviceAddr
	_, err := r.client.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("注销服务失败: %v", err)
	}
	if r.lease != 0 {
		_, _ = r.client.Revoke(ctx, r.lease)
	}
	log.Printf("[etcd] 服务注销成功: %s -> %s", serviceName, serviceAddr)
	return nil
}

func (r *EtcdRegistry) Discover(ctx context.Context, serviceName string) ([]string, error) {
	prefix := servicePrefix + serviceName + "/"
	resp, err := r.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("发现服务失败: %v", err)
	}

	var addrs []string
	for _, kv := range resp.Kvs {
		addrs = append(addrs, string(kv.Value))
	}
	return addrs, nil
}

func (r *EtcdRegistry) Watch(serviceName string, callback func(addrs []string)) {
	prefix := servicePrefix + serviceName + "/"
	ctx := context.Background()

	watchCh := r.client.Watch(ctx, prefix, clientv3.WithPrefix())

	currentAddrs, _ := r.Discover(ctx, serviceName)
	if len(currentAddrs) > 0 {
		callback(currentAddrs)
	}

	go func() {
		for range watchCh {
			addrs, err := r.Discover(ctx, serviceName)
			if err != nil {
				log.Printf("[etcd] watch 发现服务异常: %v", err)
				continue
			}
			callback(addrs)
		}
	}()
}

func (r *EtcdRegistry) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

func GetServicePrefix() string {
	return servicePrefix
}

func ParseServiceKey(key string) (serviceName, addr string, ok bool) {
	key = strings.TrimPrefix(key, servicePrefix)
	parts := strings.SplitN(key, "/", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	return parts[0], parts[1], true
}
