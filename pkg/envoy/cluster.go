package envoy

import (
	envoy_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
)

// GenerateCluster generates an Envoy cluster configuration
func GenerateCluster(name string, endpoints []string) *envoy_cluster_v3.Cluster {
	// 配置端点和其他集群设置
	// 这里只是一个示例，您需要根据实际需求来实现
	return &envoy_cluster_v3.Cluster{
		Name: name,
		// 配置其他集群属性
	}
}
