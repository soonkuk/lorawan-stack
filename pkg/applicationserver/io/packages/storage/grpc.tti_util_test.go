// Copyright © 2020 The Things Industries B.V.

package storage_test

import (
	"context"
	"fmt"
	"net"
	"time"

	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages/storage/provider"
	"go.thethings.network/lorawan-stack/v3/pkg/component"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/rpcserver"
	"go.thethings.network/lorawan-stack/v3/pkg/tenant"
	"go.thethings.network/lorawan-stack/v3/pkg/ttipb"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/unique"
	"google.golang.org/grpc/metadata"
)

type registerer struct {
	packages.ApplicationPackageHandler
}

// Roles implements the rpcserver.Registerer interface.
func (r *registerer) Roles() []ttnpb.ClusterRole {
	return nil
}

func tenantCtx(ctx context.Context, tenantID string) context.Context {
	return tenant.NewContext(ctx, ttipb.TenantIdentifiers{
		TenantID: tenantID,
	})
}

func mustHavePeer(ctx context.Context, c *component.Component, role ttnpb.ClusterRole) {
	for i := 0; i < 20; i++ {
		time.Sleep(20 * time.Millisecond)
		if _, err := c.GetPeer(ctx, role, nil); err == nil {
			return
		}
	}
	panic("could not connect to peer")
}

type mockAuth struct {
	token  string
	rights []ttnpb.Right
}

type mockIS struct {
	ttnpb.ApplicationRegistryServer
	ttnpb.ApplicationAccessServer
	applications     map[string]*ttnpb.Application
	applicationAuths map[string][]mockAuth
}

func startMockIS(ctx context.Context) (*mockIS, string) {
	is := &mockIS{
		applications:     make(map[string]*ttnpb.Application),
		applicationAuths: make(map[string][]mockAuth),
	}
	srv := rpcserver.New(ctx)
	ttnpb.RegisterApplicationRegistryServer(srv.Server, is)
	ttnpb.RegisterApplicationAccessServer(srv.Server, is)
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	go srv.Serve(lis)
	return is, lis.Addr().String()
}

func (is *mockIS) add(ctx context.Context, ids ttnpb.ApplicationIdentifiers, key string, rights ...ttnpb.Right) {
	uid := unique.ID(ctx, ids)
	is.applications[uid] = &ttnpb.Application{
		ApplicationIdentifiers: ids,
	}
	if key != "" {
		auths := is.applicationAuths[uid]
		auths = append(auths, mockAuth{
			token:  fmt.Sprintf("Bearer %v", key),
			rights: rights,
		})
		is.applicationAuths[uid] = auths
	}
}

var errNotFound = errors.DefineNotFound("not_found", "not found")

func (is *mockIS) Get(ctx context.Context, req *ttnpb.GetApplicationRequest) (*ttnpb.Application, error) {
	uid := unique.ID(ctx, req.ApplicationIdentifiers)
	app, ok := is.applications[uid]
	if !ok {
		return nil, errNotFound.New()
	}
	return app, nil
}

func (is *mockIS) ListRights(ctx context.Context, ids *ttnpb.ApplicationIdentifiers) (res *ttnpb.Rights, err error) {
	res = &ttnpb.Rights{}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return
	}
	authorization, ok := md["authorization"]
	if !ok || len(authorization) == 0 {
		return
	}
	auths, ok := is.applicationAuths[unique.ID(ctx, *ids)]
	if !ok {
		return
	}
	for _, auth := range auths {
		if auth.token == authorization[0] {
			res.Rights = append(res.Rights, auth.rights...)
		}
	}
	return
}

type mockStorage struct {
	ups []*ttnpb.ApplicationUp
}

func (p *mockStorage) Name() string { return "mock" }

func (p *mockStorage) Store(ups []*io.ContextualApplicationUp) error {
	return nil
}

func (p *mockStorage) Range(ctx context.Context, q provider.Query, f func(up *ttnpb.ApplicationUp) error) error {
	if p.ups == nil {
		return nil
	}
	for _, up := range p.ups {
		f(up)
	}
	return nil
}

func newMockStorageProvider(ups []*ttnpb.ApplicationUp) provider.Provider {
	return &mockStorage{ups}
}
