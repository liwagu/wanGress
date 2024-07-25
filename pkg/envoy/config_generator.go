package envoy

import (
	"fmt"
	"time"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	testiov1 "wanGress/api/v1"
)

func GenerateEnvoyConfig(wangress *testiov1.WanGress) ([]*listener.Listener, []*cluster.Cluster, []*route.RouteConfiguration, []*endpoint.ClusterLoadAssignment, error) {
	var listeners []*listener.Listener
	var clusters []*cluster.Cluster
	var routes []*route.RouteConfiguration
	var endpoints []*endpoint.ClusterLoadAssignment

	// Generate main route configuration
	routeConfig := &route.RouteConfiguration{
		Name: "main_route",
		VirtualHosts: []*route.VirtualHost{
			{
				Name:    "main_host",
				Domains: wangress.Spec.Hosts,
				Routes:  make([]*route.Route, len(wangress.Spec.Routes)),
			},
		},
	}

	for i, r := range wangress.Spec.Routes {
		// Generate route
		routeConfig.VirtualHosts[0].Routes[i] = &route.Route{
			Match: &route.RouteMatch{
				PathSpecifier: &route.RouteMatch_Prefix{
					Prefix: r.Path,
				},
			},
			Action: &route.Route_Route{
				Route: &route.RouteAction{
					ClusterSpecifier: &route.RouteAction_Cluster{
						Cluster: r.Services[0].Name,
					},
				},
			},
		}

		clusterConfig := &cluster.Cluster{
			Name:           r.Services[0].Name,
			ConnectTimeout: durationpb.New(5 * time.Second),
			ClusterDiscoveryType: &cluster.Cluster_Type{
				Type: cluster.Cluster_STRICT_DNS,
			},
			LbPolicy:        cluster.Cluster_ROUND_ROBIN,
			DnsLookupFamily: cluster.Cluster_V4_ONLY,
		}

		if err := clusterConfig.Validate(); err != nil {
			return nil, nil, nil, nil, fmt.Errorf("invalid cluster configuration: %v", err)
		}

		clusters = append(clusters, clusterConfig)

		// Generate endpoint
		endpoints = append(endpoints, &endpoint.ClusterLoadAssignment{
			ClusterName: r.Services[0].Name,
			Endpoints: []*endpoint.LocalityLbEndpoints{
				{
					LbEndpoints: []*endpoint.LbEndpoint{
						{
							HostIdentifier: &endpoint.LbEndpoint_Endpoint{
								Endpoint: &endpoint.Endpoint{
									Address: &core.Address{
										Address: &core.Address_SocketAddress{
											SocketAddress: &core.SocketAddress{
												Address: r.Services[0].Name,
												PortSpecifier: &core.SocketAddress_PortValue{
													PortValue: uint32(r.Services[0].Port),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		})
	}

	routes = append(routes, routeConfig)

	// Generate listener
	httpManager := &hcm.HttpConnectionManager{
		CodecType:  hcm.HttpConnectionManager_AUTO,
		StatPrefix: "ingress_http",
		RouteSpecifier: &hcm.HttpConnectionManager_Rds{
			Rds: &hcm.Rds{
				ConfigSource: &core.ConfigSource{
					ResourceApiVersion: core.ApiVersion_V3,
					ConfigSourceSpecifier: &core.ConfigSource_Ads{
						Ads: &core.AggregatedConfigSource{},
					},
				},
				RouteConfigName: "main_route",
			},
		},
		HttpFilters: []*hcm.HttpFilter{
			{
				Name: wellknown.Router,
			},
		},
	}

	pbst, err := anypb.New(httpManager)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to marshal HTTP connection manager: %v", err)
	}

	listeners = append(listeners, &listener.Listener{
		Name: "main_listener",
		Address: &core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.SocketAddress_TCP,
					Address:  "0.0.0.0",
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: 80,
					},
				},
			},
		},
		FilterChains: []*listener.FilterChain{
			{
				Filters: []*listener.Filter{
					{
						Name: wellknown.HTTPConnectionManager,
						ConfigType: &listener.Filter_TypedConfig{
							TypedConfig: pbst,
						},
					},
				},
			},
		},
	})

	return listeners, clusters, routes, endpoints, nil
}
