// Copyright © 2020 The Things Industries B.V.

package storage_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages/storage"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages/storage/mock"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages/storage/provider"
	"go.thethings.network/lorawan-stack/v3/pkg/cluster"
	"go.thethings.network/lorawan-stack/v3/pkg/component"
	componenttest "go.thethings.network/lorawan-stack/v3/pkg/component/test"
	"go.thethings.network/lorawan-stack/v3/pkg/config"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/rpcmetadata"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
	"google.golang.org/grpc"
)

func TestGRPC(t *testing.T) {
	ctx := tenantCtx(test.Context(), "foo-tenant")
	a := assertions.New(t)

	registeredApplicationID := ttnpb.ApplicationIdentifiers{
		ApplicationID: "app1",
	}
	unregisteredApplicationID := ttnpb.ApplicationIdentifiers{
		ApplicationID: "x-app1",
	}
	registeredApplicationKey := "key1"
	registeredApplicationKeyWithoutRights := "keyWithoutRights"
	registeredDeviceID := &ttnpb.EndDeviceIdentifiers{
		ApplicationIdentifiers: registeredApplicationID,
		DeviceID:               "device",
	}
	unregisteredDeviceID := &ttnpb.EndDeviceIdentifiers{
		ApplicationIdentifiers: unregisteredApplicationID,
		DeviceID:               "device",
	}

	is, isAddr := startMockIS(ctx)
	is.add(ctx, registeredApplicationID, registeredApplicationKey, ttnpb.RIGHT_APPLICATION_TRAFFIC_READ)
	is.add(ctx, registeredApplicationID, registeredApplicationKeyWithoutRights)

	c := componenttest.NewComponent(t, &component.Config{
		ServiceBase: config.ServiceBase{
			GRPC: config.GRPC{
				Listen:                      ":0",
				AllowInsecureForCredentials: true,
			},
			Cluster: cluster.Config{
				IdentityServer: isAddr,
			},
		},
	})
	p, close := mock.New(1, 1)
	defer close()

	s, err := storage.New(ctx, p)
	a.So(err, should.BeNil)

	c.RegisterGRPC(&registerer{s})
	componenttest.StartComponent(t, c)
	defer c.Close()
	mustHavePeer(ctx, c, ttnpb.ClusterRole_ENTITY_REGISTRY)

	grpcClient := ttnpb.NewApplicationUpStorageClient(c.LoopbackConn())

	t.Run("GetStoredApplicationUp", func(t *testing.T) {
		for _, tc := range []struct {
			name        string
			req         *ttnpb.GetStoredApplicationUpRequest
			ctxFunc     func() context.Context
			checkErr    func(*testing.T, error)
			callOptions func() []grpc.CallOption
			waitReq     bool
		}{
			{
				name: "NoApplicationIDs",
				req:  &ttnpb.GetStoredApplicationUpRequest{},
				ctxFunc: func() context.Context {
					return tenantCtx(test.Context(), "foo-tenant")
				},
				checkErr: func(t *testing.T, err error) {
					a := assertions.New(t)

					a.So(err, should.NotBeNil)
					a.So(errors.IsInvalidArgument(err), should.BeTrue)
				},
			},
			{
				name: "ApplicationNotFound",
				req:  (&ttnpb.GetStoredApplicationUpRequest{}).WithApplicationIDs(&unregisteredApplicationID),
				ctxFunc: func() context.Context {
					return tenantCtx(test.Context(), "foo-tenant")
				},
				checkErr: func(t *testing.T, err error) {
					a := assertions.New(t)

					a.So(err, should.NotBeNil)
					a.So(errors.IsUnauthenticated(err), should.BeTrue)
				},
			},
			{
				name: "ApplicationNoRights",
				req:  (&ttnpb.GetStoredApplicationUpRequest{}).WithApplicationIDs(&registeredApplicationID),
				ctxFunc: func() context.Context {
					return tenantCtx(test.Context(), "foo-tenant")
				},
				callOptions: func() []grpc.CallOption {
					return []grpc.CallOption{grpc.PerRPCCredentials(rpcmetadata.MD{
						AuthType:      "Bearer",
						AuthValue:     registeredApplicationKeyWithoutRights,
						AllowInsecure: true,
					})}
				},
				checkErr: func(t *testing.T, err error) {
					a := assertions.New(t)

					a.So(err, should.NotBeNil)
					a.So(errors.IsPermissionDenied(err), should.BeTrue)
				},
			},
			{
				name: "DeviceNotFound",
				req:  (&ttnpb.GetStoredApplicationUpRequest{}).WithEndDeviceIDs(unregisteredDeviceID),
				ctxFunc: func() context.Context {
					return tenantCtx(test.Context(), "foo-tenant")
				},
				checkErr: func(t *testing.T, err error) {
					a := assertions.New(t)

					a.So(err, should.NotBeNil)
					a.So(errors.IsUnauthenticated(err), should.BeTrue)
				},
			},
			{
				name: "DeviceNoRights",
				req:  (&ttnpb.GetStoredApplicationUpRequest{}).WithEndDeviceIDs(registeredDeviceID),
				ctxFunc: func() context.Context {
					return tenantCtx(test.Context(), "foo-tenant")
				},
				callOptions: func() []grpc.CallOption {
					return []grpc.CallOption{grpc.PerRPCCredentials(rpcmetadata.MD{
						AuthType:      "Bearer",
						AuthValue:     registeredApplicationKeyWithoutRights,
						AllowInsecure: true,
					})}
				},
				checkErr: func(t *testing.T, err error) {
					a := assertions.New(t)

					a.So(err, should.NotBeNil)
					a.So(errors.IsPermissionDenied(err), should.BeTrue)
				},
			},
			{
				name: "InvalidType",
				req: (&ttnpb.GetStoredApplicationUpRequest{
					Type: "invalid",
				}).WithApplicationIDs(&registeredApplicationID),
				ctxFunc: func() context.Context {
					return tenantCtx(test.Context(), "foo-tenant")
				},
				callOptions: func() []grpc.CallOption {
					return []grpc.CallOption{grpc.PerRPCCredentials(rpcmetadata.MD{
						AuthType:      "Bearer",
						AuthValue:     registeredApplicationKey,
						AllowInsecure: true,
					})}
				},
				checkErr: func(t *testing.T, err error) {
					a := assertions.New(t)
					a.So(err, should.NotBeNil)
					a.So(errors.IsInvalidArgument(err), should.BeTrue)
					a.So(errors.Attributes(err)["field"], should.Equal, "type")
				},
			},
			{
				name: "InvalidOrder",
				req: (&ttnpb.GetStoredApplicationUpRequest{
					Order: "invalid",
				}).WithApplicationIDs(&registeredApplicationID),
				ctxFunc: func() context.Context {
					return tenantCtx(test.Context(), "foo-tenant")
				},
				callOptions: func() []grpc.CallOption {
					return []grpc.CallOption{grpc.PerRPCCredentials(rpcmetadata.MD{
						AuthType:      "Bearer",
						AuthValue:     registeredApplicationKey,
						AllowInsecure: true,
					})}
				},
				checkErr: func(t *testing.T, err error) {
					a := assertions.New(t)
					a.So(err, should.NotBeNil)
					a.So(errors.IsInvalidArgument(err), should.BeTrue)
					a.So(errors.Attributes(err)["field"], should.Equal, "order")
				},
			},
			{
				name: "Application",
				req: (&ttnpb.GetStoredApplicationUpRequest{
					Order: "received_at",
					Type:  "uplink_message",
				}).WithApplicationIDs(&registeredApplicationID),
				ctxFunc: func() context.Context {
					return tenantCtx(test.Context(), "foo-tenant")
				},
				callOptions: func() []grpc.CallOption {
					return []grpc.CallOption{grpc.PerRPCCredentials(rpcmetadata.MD{
						AuthType:      "Bearer",
						AuthValue:     registeredApplicationKey,
						AllowInsecure: true,
					})}
				},
				checkErr: func(t *testing.T, err error) {
					a := assertions.New(t)
					a.So(err, should.Equal, io.EOF)
				},
				waitReq: true,
			},
			{
				name: "Device",
				req: (&ttnpb.GetStoredApplicationUpRequest{
					Order: "received_at",
					Type:  "uplink_message",
				}).WithEndDeviceIDs(registeredDeviceID),
				ctxFunc: func() context.Context {
					return tenantCtx(test.Context(), "foo-tenant")
				},
				callOptions: func() []grpc.CallOption {
					return []grpc.CallOption{grpc.PerRPCCredentials(rpcmetadata.MD{
						AuthType:      "Bearer",
						AuthValue:     registeredApplicationKey,
						AllowInsecure: true,
					})}
				},
				checkErr: func(t *testing.T, err error) {
					a := assertions.New(t)
					a.So(err, should.Equal, io.EOF)
				},
				waitReq: true,
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				a := assertions.New(t)
				opts := []grpc.CallOption{}
				if tc.callOptions != nil {
					opts = append(opts, tc.callOptions()...)
				}
				client, err := grpcClient.GetStoredApplicationUp(tc.ctxFunc(), tc.req, opts...)
				a.So(err, should.BeNil)

				_, err = client.Recv()
				tc.checkErr(t, err)

				if tc.waitReq {
					select {
					case query, ok := <-p.QueryCh:
						a.So(ok, should.BeTrue)
						a.So(query, should.Resemble, provider.QueryFromRequest(tc.req))
					case <-time.After(time.Second):
						t.Fatal("Timed out waiting for request to be fulfilled")
						t.FailNow()
					}
				}
			})
		}
	})

}
