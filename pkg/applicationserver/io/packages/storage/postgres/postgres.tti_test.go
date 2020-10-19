// Copyright © 2020 The Things Industries B.V.

package postgres_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	pbtypes "github.com/gogo/protobuf/types"
	"github.com/jinzhu/gorm"
	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages/storage/postgres"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages/storage/provider"
	"go.thethings.network/lorawan-stack/v3/pkg/random"
	"go.thethings.network/lorawan-stack/v3/pkg/tenant"
	"go.thethings.network/lorawan-stack/v3/pkg/ttipb"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/util/randutil"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
)

func tenantCtx(id string) context.Context {
	return tenant.NewContext(test.Context(), ttipb.TenantIdentifiers{
		TenantID: id,
	})
}

func timePtr(t time.Time) *time.Time { return &t }

func newUplink(r *randutil.LockedRand, ids *ttnpb.EndDeviceIdentifiers, receivedAt time.Time) *ttnpb.ApplicationUp {
	uplink := &ttnpb.ApplicationUp_UplinkMessage{
		UplinkMessage: &ttnpb.ApplicationUplink{
			FRMPayload: random.Bytes(4),
			FPort:      uint32(random.Intn(255)),
			FCnt:       uint32(random.Intn(100)),
			DecodedPayload: &pbtypes.Struct{
				Fields: map[string]*pbtypes.Value{
					"key": {
						Kind: &pbtypes.Value_StringValue{
							StringValue: "value",
						},
					},
					"integer": {
						Kind: &pbtypes.Value_NumberValue{
							NumberValue: float64(-50 + random.Intn(100)),
						},
					},
				},
			},
			ReceivedAt: receivedAt,
		},
	}
	up := &ttnpb.ApplicationUp{
		ReceivedAt:           &receivedAt,
		EndDeviceIdentifiers: *ids,
		Up:                   uplink,
	}
	return up
}

func newJoinAccept(r *randutil.LockedRand, ids *ttnpb.EndDeviceIdentifiers, receivedAt time.Time) *ttnpb.ApplicationUp {
	return &ttnpb.ApplicationUp{
		ReceivedAt:           &receivedAt,
		EndDeviceIdentifiers: *ids,
		Up:                   ttnpb.NewPopulatedApplicationUp_JoinAccept(r, true),
	}
}

func errIsNil(err error) bool { return err == nil }

func from(list []*ttnpb.ApplicationUp, indices ...int) []*ttnpb.ApplicationUp {
	res := make([]*ttnpb.ApplicationUp, 0, len(indices))
	for _, idx := range indices {
		res = append(res, list[idx])
	}
	return res
}

func TestPostgresProvider(t *testing.T) {
	withDB(t, func(t *testing.T, dbConnString string) {
		ctx := test.Context()
		p, err := postgres.New(ctx, packages.PostgresConfig{
			DatabaseURI:     dbConnString,
			Debug:           os.Getenv("SQL_DEBUG") == "1",
			SelectBatchSize: 1,
		})
		a := assertions.New(t)
		if !a.So(err, should.BeNil) {
			t.FailNow()
		}

		r := test.Randy
		now := time.Now().UTC()

		// app1 has dev1, dev3
		// app2 has dev2
		app1 := ttnpb.NewPopulatedApplicationIdentifiers(r, false)
		app2 := ttnpb.NewPopulatedApplicationIdentifiers(r, false)
		dev1 := ttnpb.NewPopulatedEndDeviceIdentifiers(r, false)
		dev2 := ttnpb.NewPopulatedEndDeviceIdentifiers(r, false)
		dev3 := ttnpb.NewPopulatedEndDeviceIdentifiers(r, false)
		dev1.ApplicationIdentifiers = *app1
		dev2.ApplicationIdentifiers = *app1
		dev3.ApplicationIdentifiers = *app2

		const (
			DEV1       = 0
			DEV1_OLDER = 1
			DEV2       = 2
			DEV3       = 3
			DEV3_JOIN  = 4
		)
		ups := []*ttnpb.ApplicationUp{
			newUplink(r, dev1, now),
			newUplink(r, dev1, now.Add(-10*time.Minute)),
			newUplink(r, dev2, now.Add(-time.Minute)),
			newUplink(r, dev3, now.Add(time.Hour)),
			newJoinAccept(r, dev3, now.Add(2*time.Hour)),
		}

		for _, up := range ups[0:4] {
			if up.GetUplinkMessage().FPort == ups[0].GetUplinkMessage().FPort {
				up.GetUplinkMessage().FPort++
			}
		}

		t.Run("Store", func(t *testing.T) {
			for _, up := range ups {
				a.So(p.Store([]*io.ContextualApplicationUp{{
					Context:       tenantCtx("tenant-1"),
					ApplicationUp: up,
				}}), should.BeNil)
			}

			for _, up := range ups {
				a.So(p.Store([]*io.ContextualApplicationUp{{
					Context:       tenantCtx("tenant-2"),
					ApplicationUp: up,
				}}), should.BeNil)
			}
		})

		t.Run("Range", func(t *testing.T) {
			for _, tc := range []struct {
				name        string
				tenantID    string
				req         *ttnpb.GetStoredApplicationUpRequest
				validateErr func(err error) bool
				expected    []*ttnpb.ApplicationUp
			}{
				{
					name:     "DeviceUplinks",
					tenantID: "tenant-1",
					req: (&ttnpb.GetStoredApplicationUpRequest{
						Type:  "uplink_message",
						Order: "received_at",
					}).WithEndDeviceIDs(dev1),
					validateErr: errIsNil,
					expected:    from(ups, DEV1_OLDER, DEV1),
				},
				{
					name:     "DeviceUplinksLimit",
					tenantID: "tenant-1",
					req: (&ttnpb.GetStoredApplicationUpRequest{
						Type:  "uplink_message",
						Order: "received_at",
						Limit: &pbtypes.UInt32Value{
							Value: 1,
						},
					}).WithEndDeviceIDs(dev1),
					validateErr: errIsNil,
					expected:    from(ups, DEV1_OLDER),
				},
				{
					name:     "DeviceOrder",
					tenantID: "tenant-1",
					req: (&ttnpb.GetStoredApplicationUpRequest{
						Type:  "uplink_message",
						Order: "-received_at",
					}).WithEndDeviceIDs(dev1),
					validateErr: errIsNil,
					expected:    from(ups, DEV1, DEV1_OLDER),
				},
				{
					name:     "DeviceOrderLimit",
					tenantID: "tenant-1",
					req: (&ttnpb.GetStoredApplicationUpRequest{
						Type:  "uplink_message",
						Order: "-received_at",
						Limit: &pbtypes.UInt32Value{
							Value: 1,
						},
					}).WithEndDeviceIDs(dev1),
					validateErr: errIsNil,
					expected:    from(ups, DEV1),
				},
				{
					name:     "DeviceType",
					tenantID: "tenant-1",
					req: (&ttnpb.GetStoredApplicationUpRequest{
						Type: "join_accept",
					}).WithEndDeviceIDs(dev1),
					validateErr: errIsNil,
					expected:    []*ttnpb.ApplicationUp{},
				},
				{
					name:     "DeviceJoinAccept",
					tenantID: "tenant-2",
					req: (&ttnpb.GetStoredApplicationUpRequest{
						Type: "join_accept",
					}).WithEndDeviceIDs(dev3),
					validateErr: errIsNil,
					expected:    from(ups, DEV3_JOIN),
				},
				{
					name:     "DeviceAllTypes",
					tenantID: "tenant-1",
					req: (&ttnpb.GetStoredApplicationUpRequest{
						Type: "",
					}).WithEndDeviceIDs(dev3),
					validateErr: errIsNil,
					expected:    from(ups, DEV3, DEV3_JOIN),
				},
				{
					name:     "DeviceOtherTenant",
					tenantID: "my-tenant",
					req: (&ttnpb.GetStoredApplicationUpRequest{
						Type: "",
					}).WithEndDeviceIDs(dev3),
					validateErr: errIsNil,
					expected:    []*ttnpb.ApplicationUp{},
				},
				{
					name:     "DeviceAfter",
					tenantID: "tenant-1",
					req: (&ttnpb.GetStoredApplicationUpRequest{
						Type:  "uplink_message",
						After: timePtr(now.Add(-time.Second)),
					}).WithEndDeviceIDs(dev1),
					validateErr: errIsNil,
					expected:    from(ups, DEV1),
				},
				{
					name:     "DeviceBefore",
					tenantID: "tenant-1",
					req: (&ttnpb.GetStoredApplicationUpRequest{
						Type:   "uplink_message",
						Before: timePtr(now.Add(-time.Second)),
					}).WithEndDeviceIDs(dev1),
					validateErr: errIsNil,
					expected:    from(ups, DEV1_OLDER),
				},
				{
					name:     "DeviceBetween",
					tenantID: "tenant-1",
					req: (&ttnpb.GetStoredApplicationUpRequest{
						EndDeviceIDs: dev1,
						Type:         "uplink_message",
						After:        timePtr(now.Add(-10 * time.Hour)),
						Before:       timePtr(now.Add(time.Hour)),
						Order:        "received_at",
					}).WithEndDeviceIDs(dev1),
					validateErr: errIsNil,
					expected:    from(ups, DEV1_OLDER, DEV1),
				},
				{
					name:     "DeviceBetweenOne",
					tenantID: "tenant-1",
					req: (&ttnpb.GetStoredApplicationUpRequest{
						Type:   "uplink_message",
						After:  timePtr(now.Add(-5 * time.Second)),
						Before: timePtr(now.Add(time.Hour)),
						Order:  "received_at",
					}).WithEndDeviceIDs(dev1),
					validateErr: errIsNil,
					expected:    from(ups, DEV1),
				},
				{
					name:     "DeviceFPort",
					tenantID: "tenant-1",
					req: (&ttnpb.GetStoredApplicationUpRequest{
						FPort: &pbtypes.UInt32Value{
							Value: ups[DEV1].GetUplinkMessage().FPort,
						},
					}).WithEndDeviceIDs(dev1),
					validateErr: errIsNil,
					expected:    from(ups, DEV1),
				},
				{
					name:     "DeviceFPortReturnOnlyUplinks",
					tenantID: "tenant-1",
					req: (&ttnpb.GetStoredApplicationUpRequest{
						FPort: &pbtypes.UInt32Value{
							Value: ups[DEV3].GetUplinkMessage().FPort,
						},
					}).WithEndDeviceIDs(dev3),
					validateErr: errIsNil,
					expected:    from(ups, DEV3),
				},
				{
					name:     "ApplicationUplinks",
					tenantID: "tenant-1",
					req: (&ttnpb.GetStoredApplicationUpRequest{
						Type:  "uplink_message",
						Order: "received_at",
					}).WithApplicationIDs(app1),
					validateErr: errIsNil,
					expected:    from(ups, DEV1_OLDER, DEV2, DEV1),
				},
				{
					name:     "ApplicationUplinksLimit",
					tenantID: "tenant-1",
					req: (&ttnpb.GetStoredApplicationUpRequest{
						Type:  "uplink_message",
						Order: "received_at",
						Limit: &pbtypes.UInt32Value{
							Value: 2,
						},
					}).WithApplicationIDs(app1),
					validateErr: errIsNil,
					expected:    from(ups, DEV1_OLDER, DEV2),
				},
				{
					name:     "ApplicationOrder",
					tenantID: "tenant-1",
					req: (&ttnpb.GetStoredApplicationUpRequest{
						Type:  "uplink_message",
						Order: "-received_at",
					}).WithApplicationIDs(app1),
					validateErr: errIsNil,
					expected:    from(ups, DEV1, DEV2, DEV1_OLDER),
				},
				{
					name:     "ApplicationOrderLimit",
					tenantID: "tenant-1",
					req: (&ttnpb.GetStoredApplicationUpRequest{
						Type:  "uplink_message",
						Order: "-received_at",
						Limit: &pbtypes.UInt32Value{
							Value: 2,
						},
					}).WithApplicationIDs(app1),
					validateErr: errIsNil,
					expected:    from(ups, DEV1, DEV2),
				},
				{
					name:     "ApplicationType",
					tenantID: "tenant-1",
					req: (&ttnpb.GetStoredApplicationUpRequest{
						Type: "join_accept",
					}).WithApplicationIDs(app2),
					validateErr: errIsNil,
					expected:    from(ups, DEV3_JOIN),
				},
				{
					name:     "ApplicationAllTypes",
					tenantID: "tenant-1",
					req: (&ttnpb.GetStoredApplicationUpRequest{
						Type: "",
					}).WithApplicationIDs(app2),
					validateErr: errIsNil,
					expected:    from(ups, DEV3, DEV3_JOIN),
				},
				{
					name:     "ApplicationOtherTenant",
					tenantID: "my-tenant",
					req: (&ttnpb.GetStoredApplicationUpRequest{
						Type: "",
					}).WithApplicationIDs(app2),
					validateErr: errIsNil,
					expected:    []*ttnpb.ApplicationUp{},
				},
				{
					name:     "ApplicationAfter",
					tenantID: "tenant-1",
					req: (&ttnpb.GetStoredApplicationUpRequest{
						Type:  "uplink_message",
						After: timePtr(now.Add(-2 * time.Minute)),
					}).WithApplicationIDs(app1),
					validateErr: errIsNil,
					expected:    from(ups, DEV1, DEV2),
				},
				{
					name:     "ApplicationBefore",
					tenantID: "tenant-1",
					req: (&ttnpb.GetStoredApplicationUpRequest{
						Type:   "uplink_message",
						Before: timePtr(now.Add(-time.Second)),
					}).WithApplicationIDs(app1),
					validateErr: errIsNil,
					expected:    from(ups, DEV1_OLDER, DEV2),
				},
				{
					name:     "ApplicationBetween",
					tenantID: "tenant-1",
					req: (&ttnpb.GetStoredApplicationUpRequest{
						Type:   "uplink_message",
						After:  timePtr(now.Add(-10 * time.Hour)),
						Before: timePtr(now.Add(time.Hour)),
						Order:  "received_at",
					}).WithApplicationIDs(app1),
					validateErr: errIsNil,
					expected:    from(ups, DEV1_OLDER, DEV2, DEV1),
				},
				{
					name:     "ApplicationBetweenOne",
					tenantID: "tenant-1",
					req: (&ttnpb.GetStoredApplicationUpRequest{
						Type:   "uplink_message",
						After:  timePtr(now.Add(-5 * time.Second)),
						Before: timePtr(now.Add(time.Hour)),
						Order:  "received_at",
					}).WithApplicationIDs(app1),
					validateErr: errIsNil,
					expected:    from(ups, DEV1),
				},
				{
					name:     "ApplicationFPort",
					tenantID: "tenant-1",
					req: (&ttnpb.GetStoredApplicationUpRequest{
						Type: "uplink_message",
						FPort: &pbtypes.UInt32Value{
							Value: ups[DEV1].GetUplinkMessage().FPort,
						},
					}).WithApplicationIDs(app1),
					validateErr: errIsNil,
					expected:    from(ups, DEV1),
				},
			} {
				t.Run(tc.name, func(t *testing.T) {
					result := []*ttnpb.ApplicationUp{}
					ctx := tenant.NewContext(test.Context(), ttipb.TenantIdentifiers{
						TenantID: tc.tenantID,
					})
					err := p.Range(ctx, provider.QueryFromRequest(tc.req), func(up *ttnpb.ApplicationUp) error {
						result = append(result, up)
						return nil
					})

					a := assertions.New(t)
					a.So(tc.validateErr(err), should.BeTrue)

					a.So(tenant.FromContext(ctx).TenantID, should.Equal, tc.tenantID)

					if !a.So(len(result), should.Resemble, len(tc.expected)) {
						t.FailNow()
					}
					a.So(result, should.Resemble, tc.expected)
				})
			}
		})
	})
}

var (
	setup   sync.Once
	setupDB *gorm.DB
)

func withDB(t *testing.T, f func(*testing.T, string)) {
	var dbConnString string
	setup.Do(func() {
		dbAddress := os.Getenv("SQL_DB_ADDRESS")
		if dbAddress == "" {
			dbAddress = "localhost:26257"
		}
		dbName := os.Getenv("TEST_DATABASE_NAME")
		if dbName == "" {
			dbName = "ttn_lorawan_is_store_test"
		}
		dbAuth := os.Getenv("SQL_DB_AUTH")
		if dbAuth == "" {
			dbAuth = "root"
		}
		dbConnString = fmt.Sprintf("postgresql://%s@%s/%s?sslmode=disable", dbAuth, dbAddress, dbName)
		var err error
		setupDB, err = postgres.Open(dbConnString, packages.PostgresConfig{
			Debug: os.Getenv("SQL_DEBUG") == "1",
		})
		if err != nil {
			panic(err)
		}
		if err := setupDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)).Error; err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				panic(err)
			}
		}
		if err = postgres.Initialize(setupDB); err != nil {
			panic(err)
		}
	})
	defer func() {
		setupDB.DropTableIfExists(postgres.StoredApplicationUp{}.TableName())
		setupDB.Close()
	}()
	f(t, dbConnString)
}
