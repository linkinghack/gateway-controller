package xds

import (
	"context"
	"net"
	"time"

	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/linkinghack/gateway-controller/config"
	"github.com/linkinghack/gateway-controller/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	clusterservice "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoverygrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpointservice "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	listenerservice "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	routeservice "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	runtimeservice "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	secretservice "github.com/envoyproxy/go-control-plane/envoy/service/secret/v3"
)

type XdsServer struct {
	cache       cache.SnapshotCache
	snapVersion uint64
	envoyServer server.Server
	grpcServer  *grpc.Server

	ListenAddress string // ip:port
}

func NewXdsServerWithGlobalConfig(ctx context.Context) *XdsServer {
	log := log.GetSpecificLogger("XdsServer")
	xdsConf := config.GetGlobalConfig().XdsServerConfig

	// Create XDS server config cache
	c := cache.NewSnapshotCache(false, cache.IDHash{}, log)

	// Prepare grpc server
	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions,
		grpc.MaxConcurrentStreams(uint32(xdsConf.GrpcMaxConcurrentStreams)),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    time.Duration(xdsConf.GrpcKeepaliveSeconds) * time.Second,
			Timeout: time.Duration(xdsConf.GrpcKeepaliveTimeoutSeconds) * time.Second,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             time.Duration(xdsConf.GrpcKeepaliveMinTimeSeconds) * time.Second,
			PermitWithoutStream: true,
		}),
	)
	grpcServer := grpc.NewServer(grpcOptions...)

	// Prepare Envoy XDS server
	srv := server.NewServer(ctx, c, nil)

	// register XDS services to GRPC server
	registerServer(grpcServer, srv)

	return &XdsServer{
		cache:         c,
		snapVersion:   0,
		envoyServer:   srv,
		grpcServer:    grpcServer,
		ListenAddress: xdsConf.ListenAddr,
	}
}

func (x *XdsServer) Start() {
	log := log.GetSpecificLogger("XdsServer.Start()")
	lis, err := net.Listen("tcp", x.ListenAddress)
	if err != nil {
		log.WithError(err).Fatal("Cannot create listener for XDS server")
		return
	}

	// start to listen
	if err := x.grpcServer.Serve(lis); err != nil {
		log.WithError(err).Error("GRPC server stopped")
	}
}

func (x *XdsServer) Stop() {
	log := log.GetSpecificLogger("XdsServer.Stop()")
	log.Info("Stopping XDS server")
	x.grpcServer.GracefulStop()
}

func registerServer(grpcServer *grpc.Server, server server.Server) {
	// register services
	discoverygrpc.RegisterAggregatedDiscoveryServiceServer(grpcServer, server)
	endpointservice.RegisterEndpointDiscoveryServiceServer(grpcServer, server)
	clusterservice.RegisterClusterDiscoveryServiceServer(grpcServer, server)
	routeservice.RegisterRouteDiscoveryServiceServer(grpcServer, server)
	listenerservice.RegisterListenerDiscoveryServiceServer(grpcServer, server)
	secretservice.RegisterSecretDiscoveryServiceServer(grpcServer, server)
	runtimeservice.RegisterRuntimeDiscoveryServiceServer(grpcServer, server)
}

func (x *XdsServer) GetXdsCache() cache.SnapshotCache {
	return x.cache
}
