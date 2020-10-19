// Copyright © 2020 The Things Industries B.V.

package storage_test

import (
	"testing"
	"time"

	"github.com/mohae/deepcopy"
	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages/storage"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages/storage/mock"
	"go.thethings.network/lorawan-stack/v3/pkg/events"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/util/randutil"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
)

const (
	timeout = (1 << 5) * time.Millisecond
)

func newUplink(r *randutil.LockedRand) *ttnpb.ApplicationUp {
	up := ttnpb.NewPopulatedApplicationUp(r, false)
	uplink := ttnpb.NewPopulatedApplicationUp_UplinkMessage(r, false)
	uplink.UplinkMessage.SessionKeyID = []byte("session")
	uplink.UplinkMessage.RxMetadata = append(uplink.UplinkMessage.RxMetadata, &ttnpb.RxMetadata{
		UplinkToken: []byte("uplink_token"),
	})
	up.Up = uplink
	return up
}

func TestStorage(t *testing.T) {
	a := assertions.New(t)
	r := test.Randy

	ctx := test.Context()
	p, close := mock.New(1, 1)
	defer close()

	pkg, err := storage.New(ctx, p, storage.WithEnabled(packages.StorageTypesConfig{
		All: true,
	}))
	a.So(err, should.BeNil)

	a.So(pkg.Package().Name, should.Equal, "storage-integration")

	eventsCh := map[string]events.Channel{}
	for _, event := range []string{"as.packages.storage.up.store"} {
		eventsCh[event] = make(events.Channel, 1)
	}
	defer test.SetDefaultEventsPubSub(&test.MockEventPubSub{
		PublishFunc: func(ev events.Event) {
			switch name := ev.Name(); name {
			case "as.packages.storage.up.store":
				go func() {
					eventsCh[name] <- ev
				}()
			default:
				t.Logf("%s event published", name)
			}
		},
	})()

	t.Run("HandleUp/NoAssociation", func(t *testing.T) {
		a.So(pkg.HandleUp(ctx, nil, nil, nil), should.BeNil)

		select {
		case <-p.StoreCh:
			t.Fatal("Unexpected store request")
			t.FailNow()
		case <-time.After(100 * time.Millisecond):
		}
	})

	t.Run("HandleUp/Assoc", func(t *testing.T) {
		up := newUplink(r)
		original := deepcopy.Copy(up).(*ttnpb.ApplicationUp)
		a.So(pkg.HandleUp(ctx, nil, &ttnpb.ApplicationPackageAssociation{}, up), should.BeNil)

		select {
		case recv, ok := <-p.StoreCh:
			original = storage.CleanApplicationUp(original)
			a.So(ok, should.BeTrue)
			a.So(recv, should.Resemble, []*io.ContextualApplicationUp{{
				Context:       ctx,
				ApplicationUp: original,
			}})
		case <-time.After(time.Second):
			t.Fatal("Timed out waiting for upstream messages")
			t.FailNow()
		}

		select {
		case <-eventsCh["as.packages.storage.up.store"]:
		case <-time.After(time.Second):
			t.Fatal("Timed out waiting for store event")
			t.FailNow()
		}
	})

	t.Run("HandleUp/DefaultAssoc", func(t *testing.T) {
		up := newUplink(r)
		original := deepcopy.Copy(up).(*ttnpb.ApplicationUp)
		a.So(pkg.HandleUp(ctx, &ttnpb.ApplicationPackageDefaultAssociation{}, nil, up), should.BeNil)

		select {
		case recv, ok := <-p.StoreCh:
			original = storage.CleanApplicationUp(original)
			a.So(ok, should.BeTrue)
			a.So(recv, should.Resemble, []*io.ContextualApplicationUp{{
				Context:       ctx,
				ApplicationUp: original,
			}})
		case <-time.After(time.Second):
			t.Fatal("Timed out waiting for upstream messages")
			t.FailNow()
		}

		select {
		case <-eventsCh["as.packages.storage.up.store"]:
		case <-time.After(time.Second):
			t.Fatal("Timed out waiting for store event")
			t.FailNow()
		}
	})
	t.Run("HandleUp/BothAssoc", func(t *testing.T) {
		up := newUplink(r)
		original := deepcopy.Copy(up).(*ttnpb.ApplicationUp)
		a.So(pkg.HandleUp(ctx, &ttnpb.ApplicationPackageDefaultAssociation{}, &ttnpb.ApplicationPackageAssociation{}, up), should.BeNil)

		select {
		case recv, ok := <-p.StoreCh:
			original = storage.CleanApplicationUp(original)
			a.So(ok, should.BeTrue)
			a.So(recv, should.Resemble, []*io.ContextualApplicationUp{{
				Context:       ctx,
				ApplicationUp: original,
			}})
		case <-time.After(time.Second):
			t.Fatal("Timed out waiting for upstream messages")
			t.FailNow()
		}

		select {
		case <-eventsCh["as.packages.storage.up.store"]:
		case <-time.After(time.Second):
			t.Fatal("Timed out waiting for store event")
			t.FailNow()
		}
	})

	t.Run("HandleUp/Type", func(t *testing.T) {
		for _, tc := range []struct {
			name   string
			setup  func(r *randutil.LockedRand, up *ttnpb.ApplicationUp)
			config packages.StorageTypesConfig
		}{
			{
				name: "UplinkMessage",
				setup: func(r *randutil.LockedRand, up *ttnpb.ApplicationUp) {
					up.Up = ttnpb.NewPopulatedApplicationUp_UplinkMessage(r, true)
				},
				config: packages.StorageTypesConfig{
					UplinkMessage: true,
				},
			},
			{
				name: "JoinAccept",
				setup: func(r *randutil.LockedRand, up *ttnpb.ApplicationUp) {
					up.Up = ttnpb.NewPopulatedApplicationUp_JoinAccept(r, true)
				},
				config: packages.StorageTypesConfig{
					JoinAccept: true,
				},
			},
			{
				name: "DownlinkAck",
				setup: func(r *randutil.LockedRand, up *ttnpb.ApplicationUp) {
					up.Up = ttnpb.NewPopulatedApplicationUp_DownlinkAck(r, true)
				},
				config: packages.StorageTypesConfig{
					DownlinkAck: true,
				},
			},
			{
				name: "DownlinkNack",
				setup: func(r *randutil.LockedRand, up *ttnpb.ApplicationUp) {
					up.Up = ttnpb.NewPopulatedApplicationUp_DownlinkNack(r, true)
				},
				config: packages.StorageTypesConfig{
					DownlinkNack: true,
				},
			},
			{
				name: "DownlinkSent",
				setup: func(r *randutil.LockedRand, up *ttnpb.ApplicationUp) {
					up.Up = ttnpb.NewPopulatedApplicationUp_DownlinkSent(r, true)
				},
				config: packages.StorageTypesConfig{
					DownlinkSent: true,
				},
			},
			{
				name: "DownlinkFailed",
				setup: func(r *randutil.LockedRand, up *ttnpb.ApplicationUp) {
					up.Up = ttnpb.NewPopulatedApplicationUp_DownlinkFailed(r, true)
				},
				config: packages.StorageTypesConfig{
					DownlinkFailed: true,
				},
			},
			{
				name: "DownlinkQueued",
				setup: func(r *randutil.LockedRand, up *ttnpb.ApplicationUp) {
					up.Up = ttnpb.NewPopulatedApplicationUp_DownlinkQueued(r, true)
				},
				config: packages.StorageTypesConfig{
					DownlinkQueued: true,
				},
			},
			{
				name: "DownlinkQueueInvalidated",
				setup: func(r *randutil.LockedRand, up *ttnpb.ApplicationUp) {
					up.Up = ttnpb.NewPopulatedApplicationUp_DownlinkQueueInvalidated(r, true)
				},
				config: packages.StorageTypesConfig{
					DownlinkQueueInvalidated: true,
				},
			},
			{
				name: "LocationSolved",
				setup: func(r *randutil.LockedRand, up *ttnpb.ApplicationUp) {
					up.Up = ttnpb.NewPopulatedApplicationUp_LocationSolved(r, true)
				},
				config: packages.StorageTypesConfig{
					LocationSolved: true,
				},
			},
			{
				name: "ServiceData",
				setup: func(r *randutil.LockedRand, up *ttnpb.ApplicationUp) {
					up.Up = ttnpb.NewPopulatedApplicationUp_ServiceData(r, true)
				},
				config: packages.StorageTypesConfig{
					ServiceData: true,
				},
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				a := assertions.New(t)
				r := test.Randy
				up := ttnpb.NewPopulatedApplicationUp(r, true)
				tc.setup(r, up)

				t.Run("Disabled", func(t *testing.T) {
					pkg, err := storage.New(ctx, p)
					a := assertions.New(t)
					a.So(err, should.BeNil)
					a.So(pkg.HandleUp(ctx, &ttnpb.ApplicationPackageDefaultAssociation{}, nil, up), should.BeNil)

					select {
					case <-p.StoreCh:
						t.Fatal("Provider received upstream message when none were expected")
						t.FailNow()
					case <-time.After(timeout):
					}
				})

				for _, ttc := range []struct {
					name string
					cfg  packages.StorageTypesConfig
				}{
					{
						name: "All",
						cfg:  packages.StorageTypesConfig{All: true},
					},
					{
						name: "One",
						cfg:  tc.config,
					},
				} {
					t.Run(ttc.name, func(t *testing.T) {
						pkg, err := storage.New(ctx, p, storage.WithEnabled(ttc.cfg))
						a.So(err, should.BeNil)

						a.So(pkg.HandleUp(ctx, &ttnpb.ApplicationPackageDefaultAssociation{}, nil, up), should.BeNil)

						select {
						case recv, ok := <-p.StoreCh:
							a.So(ok, should.BeTrue)
							a.So(recv, should.NotBeNil)
						case <-time.After(time.Second):
							t.Fatal("Timed out waiting for upstream messages")
							t.FailNow()
						}

						select {
						case <-eventsCh["as.packages.storage.up.store"]:
						case <-time.After(time.Second):
							t.Fatal("Timed out waiting for store event")
							t.FailNow()
						}
					})
				}
			})
		}
	})
}
