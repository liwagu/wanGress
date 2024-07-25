package envoy

import (
	"context"
	"fmt"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	clusterservice "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoverygrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpointservice "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	listenerservice "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	routeservice "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"log"
	"net"
	"time"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"

	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	xds "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"google.golang.org/grpc"
)

const (
	xdsPort = 18000
)

type XDSClient struct {
	cache  cache.SnapshotCache
	server xds.Server
}

func NewXDSClient() *XDSClient {
	cache := cache.NewSnapshotCache(false, cache.IDHash{}, nil)
	server := xds.NewServer(context.Background(), cache, nil)
	return &XDSClient{
		cache:  cache,
		server: server,
	}
}

func (c *XDSClient) Run() {
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", xdsPort))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		grpcServer := grpc.NewServer()

		// Register all xDS services
		discoverygrpc.RegisterAggregatedDiscoveryServiceServer(grpcServer, c.server)
		endpointservice.RegisterEndpointDiscoveryServiceServer(grpcServer, c.server)
		clusterservice.RegisterClusterDiscoveryServiceServer(grpcServer, c.server)
		routeservice.RegisterRouteDiscoveryServiceServer(grpcServer, c.server)
		listenerservice.RegisterListenerDiscoveryServiceServer(grpcServer, c.server)

		log.Printf("xDS server listening on %d", xdsPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}

func (c *XDSClient) UpdateConfig(listeners []*listener.Listener, clusters []*cluster.Cluster, routes []*route.RouteConfiguration, endpoints []*endpoint.ClusterLoadAssignment) error {
	version := fmt.Sprintf("%d", time.Now().UnixNano())

	snapshot, err := cache.NewSnapshot(version,
		map[resource.Type][]types.Resource{
			resource.ListenerType: castToResource(listeners),
			resource.ClusterType:  castToResource(clusters),
			resource.RouteType:    castToResource(routes),
			resource.EndpointType: castToResource(endpoints),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create snapshot: %v", err)
	}

	if err := c.cache.SetSnapshot(context.Background(), "envoy-node", snapshot); err != nil {
		return fmt.Errorf("failed to set snapshot: %v", err)
	}

	return nil
}

func (c *XDSClient) RemoveConfig() error {
	emptySnapshot, err := cache.NewSnapshot("empty",
		map[resource.Type][]types.Resource{
			resource.ListenerType: {},
			resource.ClusterType:  {},
			resource.RouteType:    {},
			resource.EndpointType: {},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create empty snapshot: %v", err)
	}

	if err := c.cache.SetSnapshot(context.Background(), "envoy-node", emptySnapshot); err != nil {
		return fmt.Errorf("failed to set empty snapshot: %v", err)
	}

	return nil
}

// Helper function to cast slices to []types.Resource
func castToResource(slice interface{}) []types.Resource {
	switch s := slice.(type) {
	case []*listener.Listener:
		r := make([]types.Resource, len(s))
		for i, v := range s {
			r[i] = v
		}
		return r
	case []*cluster.Cluster:
		r := make([]types.Resource, len(s))
		for i, v := range s {
			r[i] = v
		}
		return r
	case []*route.RouteConfiguration:
		r := make([]types.Resource, len(s))
		for i, v := range s {
			r[i] = v
		}
		return r
	default:
		return nil
	}
}
