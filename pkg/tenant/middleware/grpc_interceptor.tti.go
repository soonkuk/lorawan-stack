// Copyright © 2019 The Things Industries B.V.

package middleware

import (
	"context"
	"strings"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.thethings.network/lorawan-stack/v3/pkg/license"
	"go.thethings.network/lorawan-stack/v3/pkg/tenant"
	"go.thethings.network/lorawan-stack/v3/pkg/ttipb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func fromRPCContext(ctx context.Context, config tenant.Config) ttipb.TenantIdentifiers {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if id, ok := md["tenant-id"]; ok {
			return ttipb.TenantIdentifiers{TenantID: id[0]}
		}
		if host, ok := md["x-forwarded-host"]; ok { // Set by gRPC gateway.
			return ttipb.TenantIdentifiers{TenantID: tenantID(host[0], config)}
		}
		if authority, ok := md[":authority"]; ok { // Set by gRPC clients.
			if authority[0] != "in-process" {
				return ttipb.TenantIdentifiers{TenantID: tenantID(authority[0], config)}
			}
		}
	}
	return ttipb.TenantIdentifiers{}
}

// UnaryClientInterceptor is a gRPC interceptor that injects the tenant ID into the metadata.
func UnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	if license.RequireMultiTenancy(ctx) == nil {
		if tenantID := tenant.FromContext(ctx); !tenantID.IsZero() {
			md, _ := metadata.FromOutgoingContext(ctx)
			ctx = metadata.NewOutgoingContext(ctx, metadata.Join(md, metadata.Pairs("tenant-id", tenantID.TenantID)))
		}
	}
	return invoker(ctx, method, req, reply, cc, opts...)
}

// StreamClientInterceptor is a gRPC interceptor that injects the tenant ID into the metadata.
func StreamClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	if license.RequireMultiTenancy(ctx) == nil {
		if tenantID := tenant.FromContext(ctx); !tenantID.IsZero() {
			md, _ := metadata.FromOutgoingContext(ctx)
			ctx = metadata.NewOutgoingContext(ctx, metadata.Join(md, metadata.Pairs("tenant-id", tenantID.TenantID)))
		}
	}
	return streamer(ctx, desc, cc, method, opts...)
}

const ctxTagName = "grpc.context.tenant_id"

func extractFromRPC(ctx context.Context, config tenant.Config) context.Context {
	if license.RequireMultiTenancy(ctx) == nil {
		if id := tenant.FromContext(ctx); !id.IsZero() {
			grpc_ctxtags.Extract(ctx).Set(ctxTagName, id.TenantID)
			return ctx
		}
		if id := fromRPCContext(ctx, config); !id.IsZero() {
			grpc_ctxtags.Extract(ctx).Set(ctxTagName, id.TenantID)
			ctx = tenant.NewContext(ctx, id)
			return ctx
		}
	}
	if id := config.DefaultID; id != "" {
		grpc_ctxtags.Extract(ctx).Set(ctxTagName, id)
		ctx = tenant.NewContext(ctx, ttipb.TenantIdentifiers{TenantID: id})
		return ctx
	}
	return ctx
}

// UnaryServerExtractor returns an interceptor that extracts the tenant ID from unary RPCs.
func UnaryServerExtractor(config tenant.Config) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(extractFromRPC(ctx, config), req)
	}
}

// StreamServerExtractor returns an interceptor that extracts the tenant ID from streaming RPCs.
func StreamServerExtractor(config tenant.Config) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = extractFromRPC(stream.Context(), config)
		return handler(srv, wrapped)
	}
}

var tenantAgnosticServices = []string{"/tti.lorawan.v3.TenantRegistry", "/tti.lorawan.v3.Tbs"}

// UnaryServerFetchInterceptor returns an interceptor that fetches the tenant if there is a multi-tenant license.
func UnaryServerFetchInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		id := tenant.FromContext(ctx)
		if id.IsZero() {
			for _, service := range tenantAgnosticServices {
				if strings.HasPrefix(info.FullMethod, service) {
					return handler(ctx, req)
				}
			}
			return nil, errMissingTenantID.New()
		}
		if license.RequireMultiTenancy(ctx) == nil {
			if err := fetchTenant(ctx); err != nil {
				return nil, err
			}
		}
		return handler(ctx, req)
	}
}

// StreamServerFetchInterceptor returns an interceptor that fetches the tenant if there is a multi-tenant license.
func StreamServerFetchInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()
		id := tenant.FromContext(ctx)
		if id.IsZero() {
			for _, service := range tenantAgnosticServices {
				if strings.HasPrefix(info.FullMethod, service) {
					return handler(srv, stream)
				}
			}
			return errMissingTenantID.New()
		}
		if license.RequireMultiTenancy(ctx) == nil {
			if err := fetchTenant(ctx); err != nil {
				return err
			}
			return handler(srv, stream)
		}
		return handler(srv, stream)
	}
}
