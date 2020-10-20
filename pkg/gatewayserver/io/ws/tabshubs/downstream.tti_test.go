// Copyright © 2020 The Things Industries B.V.

package tabshubs

import (
	"context"
	"testing"
	"time"

	"go.thethings.network/lorawan-stack/v3/pkg/tenant"
	"go.thethings.network/lorawan-stack/v3/pkg/ttipb"

	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/gatewayserver/io/ws"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/unique"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
)

func timePtr(time time.Time) *time.Time { return &time }

func TestFromDownlinkMessage(t *testing.T) {
	var th tabsHubs
	baseCtx := context.Background()
	ctx := tenant.NewContext(baseCtx, ttipb.TenantIdentifiers{
		TenantID: "tti",
	})
	uid := unique.ID(ctx, ttnpb.GatewayIdentifiers{GatewayID: "test-gateway"})
	sessionCtx := ws.NewContextWithSession(ctx, &ws.Session{
		Data: State{
			ID: 0x11,
		},
	})
	for _, tc := range []struct {
		Name                    string
		DownlinkMessage         ttnpb.DownlinkMessage
		ExpectedDownlinkMessage DownlinkMessage
	}{
		{
			Name: "SampleDownlink",
			DownlinkMessage: ttnpb.DownlinkMessage{
				RawPayload: []byte("Ymxhamthc25kJ3M=="),
				EndDeviceIDs: &ttnpb.EndDeviceIdentifiers{
					DeviceID: "testdevice",
				},
				Settings: &ttnpb.DownlinkMessage_Scheduled{
					Scheduled: &ttnpb.TxSettings{
						DataRateIndex: 2,
						Frequency:     868500000,
						Downlink: &ttnpb.TxSettings_Downlink{
							AntennaIndex: 2,
						},
						Timestamp: 1553300787,
					},
				},
				CorrelationIDs: []string{"correlation1"},
			},
			ExpectedDownlinkMessage: DownlinkMessage{
				DevEUI:  "00-00-00-00-00-00-00-00",
				SeqNo:   1,
				Pdu:     "596d7868616d74686332356b4a334d3d3d",
				MuxTime: 1554300787.123456,
				Freq:    868500000,
				DR:      2,
			},
		},
		{
			Name: "WithAbsoluteTime",
			DownlinkMessage: ttnpb.DownlinkMessage{
				RawPayload: []byte("Ymxhamthc25kJ3M=="),
				EndDeviceIDs: &ttnpb.EndDeviceIdentifiers{
					DeviceID: "testdevice",
				},
				Settings: &ttnpb.DownlinkMessage_Scheduled{
					Scheduled: &ttnpb.TxSettings{
						DataRateIndex: 2,
						Frequency:     869525000,
						Downlink: &ttnpb.TxSettings_Downlink{
							AntennaIndex: 2,
						},
					},
				},
				CorrelationIDs: []string{"correlation2"},
			},
			ExpectedDownlinkMessage: DownlinkMessage{
				DevEUI:  "00-00-00-00-00-00-00-00",
				SeqNo:   2,
				Pdu:     "596d7868616d74686332356b4a334d3d3d",
				Freq:    869525000,
				DR:      2,
				MuxTime: 1554300787.123456,
			},
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			a := assertions.New(t)
			raw, err := th.FromDownlink(sessionCtx, uid, tc.DownlinkMessage, 1554300787, time.Unix(1554300787, 123456000))
			a.So(err, should.BeNil)
			var dnmsg DownlinkMessage
			err = dnmsg.unmarshalJSON(raw)
			a.So(err, should.BeNil)
			dnmsg.XTime = tc.ExpectedDownlinkMessage.XTime
			if !a.So(dnmsg, should.Resemble, tc.ExpectedDownlinkMessage) {
				t.Fatalf("Invalid DownlinkMessage: %v", dnmsg)
			}
		})
	}
}

func TestToDownlinkMessage(t *testing.T) {
	for _, tc := range []struct {
		Name                    string
		DownlinkMessage         DownlinkMessage
		ExpectedDownlinkMessage ttnpb.DownlinkMessage
	}{
		{
			Name: "SampleDownlink",
			DownlinkMessage: DownlinkMessage{
				Pdu:   "Ymxhamthc25kJ3M==",
				XTime: 1554300785,
				Freq:  868500000,
				DR:    2,
			},
			ExpectedDownlinkMessage: ttnpb.DownlinkMessage{
				RawPayload: []byte("Ymxhamthc25kJ3M=="),
				Settings: &ttnpb.DownlinkMessage_Scheduled{
					Scheduled: &ttnpb.TxSettings{
						DataRateIndex: 2,
						Frequency:     868500000,
						Timestamp:     1554300785,
					},
				},
			},
		},
		{
			Name: "WithAbsoluteTime",
			DownlinkMessage: DownlinkMessage{
				Pdu:  "Ymxhamthc25kJ3M==",
				Freq: 868500000,
				DR:   2,
			},
			ExpectedDownlinkMessage: ttnpb.DownlinkMessage{
				RawPayload: []byte("Ymxhamthc25kJ3M=="),
				Settings: &ttnpb.DownlinkMessage_Scheduled{
					Scheduled: &ttnpb.TxSettings{
						DataRateIndex: 2,
						Frequency:     868500000,
					},
				},
			},
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			a := assertions.New(t)
			dlMesg := tc.DownlinkMessage.ToDownlinkMessage()
			if !a.So(dlMesg, should.Resemble, tc.ExpectedDownlinkMessage) {
				t.Fatalf("Invalid DownlinkMessage: %v", dlMesg)
			}
		})
	}
}
