package envoy

import (
	"context"
	"fmt"
	envoy_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"google.golang.org/grpc"
	"log"
	"net"
)

const xdsPort = 8080 // 定义xDS服务器的端口

type XDSClient struct {
	server server.Server
	cache  cache.SnapshotCache
}

func NewXDSClient() *XDSClient {
	cache := cache.NewSnapshotCache(false, cache.IDHash{}, nil)
	srv := server.NewServer(context.Background(), cache, nil)

	return &XDSClient{
		server: srv,
		cache:  cache,
	}
}

func (c *XDSClient) Run() {
	// 启动xDS服务器
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", xdsPort))
		if err != nil {
			log.Fatalf("failed to listen on port %d: %v", xdsPort, err)
		}
		grpcServer := grpc.NewServer()
		server.RegisterServer(grpcServer, c.server)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC server: %v", err)
		}
	}()
}

func (c *XDSClient) UpdateConfig(listener *envoy_listener_v3.Listener, clusters []*envoy_cluster_v3.Cluster) error {
	// 创建并应用新的配置快照
	version := "version_1" // 可以根据需要动态生成版本号

	snapshot := cache.NewSnapshot(
		version,
		nil,                        // endpoints
		[]types.Resource{listener}, // 使用从generateEnvoyConfig返回的listener
		clusters,                   // 使用生成的集群配置
		nil,                        // routes
		nil,                        // runtimes
		nil,                        // secrets
	)

	// 应用快照到xDS服务器的缓存中
	if err := c.cache.SetSnapshot(context.Background(), "default", snapshot); err != nil {
		return fmt.Errorf("error setting snapshot: %w", err)
	}

	return nil
}

func (c *XDSClient) RemoveConfig() error {
	// 使用空快照来删除配置
	version := "version_empty"

	snapshot := cache.NewSnapshot(
		version,
		nil, // endpoints
		nil, // listeners
		nil, // clusters
		nil, // routes
		nil, // runtimes
		nil, // secrets
	)

	// 应用空快照到xDS服务器的缓存中
	if err := c.cache.SetSnapshot(context.Background(), "default", snapshot); err != nil {
		return fmt.Errorf("error removing snapshot: %w", err)
	}

	return nil
}
