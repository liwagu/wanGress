package envoy

import (
	envoy_api_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_http_connection_manager_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"google.golang.org/protobuf/types/known/anypb"
)

// GenerateListener generates an Envoy listener configuration
func GenerateListener(name, address string, port uint32) (*envoy_listener_v3.Listener, error) {
	manager := &envoy_http_connection_manager_v3.HttpConnectionManager{
		CodecType:  envoy_http_connection_manager_v3.HttpConnectionManager_AUTO,
		StatPrefix: "ingress_http",
		// 配置其他HTTP连接管理器属性
	}

	pbst, err := anypb.New(manager)
	if err != nil {
		return nil, err
	}

	listener := &envoy_listener_v3.Listener{
		Name: name,
		Address: &envoy_api_v3.Address{
			Address: &envoy_api_v3.Address_SocketAddress{
				SocketAddress: &envoy_api_v3.SocketAddress{
					Protocol: envoy_api_v3.SocketAddress_TCP,
					Address:  address,
					PortSpecifier: &envoy_api_v3.SocketAddress_PortValue{
						PortValue: port,
					},
				},
			},
		},
		FilterChains: []*envoy_listener_v3.FilterChain{
			{
				Filters: []*envoy_listener_v3.Filter{
					{
						Name: "envoy.filters.network.http_connection_manager",
						ConfigType: &envoy_listener_v3.Filter_TypedConfig{
							TypedConfig: pbst,
						},
					},
				},
			},
		},
	}

	return listener, nil
}
