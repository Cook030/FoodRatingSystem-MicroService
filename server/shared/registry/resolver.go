package registry

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

const scheme = "etcd"

type EtcdBuilder struct {
	Endpoints []string
}

type etcdResolver struct {
	client      *clientv3.Client
	serviceName string
	cc          resolver.ClientConn
	ctx         context.Context
	cancel      context.CancelFunc
}

func RegisterEtcdResolver(endpoints []string) {
	resolver.Register(&EtcdBuilder{Endpoints: endpoints})
}

func NewEtcdGrpcConn(etcdEndpoints []string, serviceName string) (*grpc.ClientConn, error) {
	RegisterEtcdResolver(etcdEndpoints)

	conn, err := grpc.NewClient(
		fmt.Sprintf("etcd:///%s", serviceName),
		grpc.WithResolvers(&EtcdBuilder{Endpoints: etcdEndpoints}),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingConfig": [{"%s": {}}]}`, roundrobin.Name)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("创建 gRPC 连接失败 (%s): %v", serviceName, err)
	}
	return conn, nil
}

func (b *EtcdBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   b.Endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	r := &etcdResolver{
		client:      cli,
		serviceName: target.URL.Host,
		cc:          cc,
		ctx:         ctx,
		cancel:      cancel,
	}

	go r.watch()

	return r, nil
}

func (b *EtcdBuilder) Scheme() string {
	return scheme
}

func (r *etcdResolver) watch() {
	prefix := servicePrefix + r.serviceName + "/"

	addrs, err := r.resolve()
	if err == nil && len(addrs) > 0 {
		r.updateState(addrs)
	}

	watchCh := r.client.Watch(r.ctx, prefix, clientv3.WithPrefix())
	for {
		select {
		case <-r.ctx.Done():
			return
		case watchResp, ok := <-watchCh:
			if !ok {
				return
			}
			if err := watchResp.Err(); err != nil {
				log.Printf("[etcd resolver] watch 错误: %v", err)
				continue
			}
			addrs, err := r.resolve()
			if err != nil {
				log.Printf("[etcd resolver] 解析服务失败: %v", err)
				continue
			}
			r.updateState(addrs)
		}
	}
}

func (r *etcdResolver) resolve() ([]resolver.Address, error) {
	prefix := servicePrefix + r.serviceName + "/"
	resp, err := r.client.Get(r.ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var addrs []resolver.Address
	for _, kv := range resp.Kvs {
		addrs = append(addrs, resolver.Address{Addr: string(kv.Value)})
	}
	return addrs, nil
}

func (r *etcdResolver) updateState(addrs []resolver.Address) {
	if len(addrs) == 0 {
		return
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
	log.Printf("[etcd resolver] 服务 %s 发现 %d 个实例", r.serviceName, len(addrs))
}

func (r *etcdResolver) ResolveNow(opts resolver.ResolveNowOptions) {}

func (r *etcdResolver) Close() {
	r.cancel()
	if r.client != nil {
		r.client.Close()
	}
}
