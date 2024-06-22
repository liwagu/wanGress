/*
 * @Author: liwa guliwa@foxmail.com
 * @Date: 2024-03-21 16:55:46
 * @LastEditors: liwa guliwa@foxmail.com
 * @LastEditTime: 2024-03-27 08:21:02
 * @FilePath: \wanGress\internal\controller\xds_client.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package controller

import (
	"context"
	"fmt"
	envoy_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"google.golang.org/grpc"
	"net"
)

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
		lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", xdsPort))
		grpcServer := grpc.NewServer()
		server.RegisterServer(grpcServer, c.server)
		if err := grpcServer.Serve(lis); err != nil {
			// 处理错误
		}
	}()
}

func (c *XDSClient) UpdateConfig(listener *envoy_listener_v3.Listener) error {
	// 创建并应用新的配置快照
	// 这里是一个示例，您需要根据实际的Envoy配置来创建快照
	snapshot := cache.NewSnapshot(
		"version_1",
		nil, // endpoints
		map[resource.Type][]types.Resource{
			resource.ListenerType: {listener}, // 使用从generateEnvoyConfig返回的listener
		},
		nil, // routes
		nil, // runtimes
		nil, // secrets
	)

	// 应用快照到xDS服务器的缓存中
	if err := c.cache.SetSnapshot(context.Background(), "default", snapshot); err != nil {
		return fmt.Errorf("error setting snapshot: %w", err)
	}

	return nil
}
