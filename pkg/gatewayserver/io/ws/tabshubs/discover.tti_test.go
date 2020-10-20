// Copyright © 2020 The Things Industries B.V.

package tabshubs

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/basicstation"
	"go.thethings.network/lorawan-stack/v3/pkg/gatewayserver/io/ws"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/types"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
)

func TestDiscover(t *testing.T) {
	a := assertions.New(t)
	ctx := context.Background()
	var th tabsHubs
	eui := types.EUI64{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}
	mockServer := mockServer{
		ids: ttnpb.GatewayIdentifiers{
			GatewayID: "eui-1111111111111111",
			EUI:       &eui,
		},
	}
	info := ws.ServerInfo{
		Scheme:  "wss",
		Address: "thethings.example.com:8887",
	}

	for _, tc := range []struct {
		Name             string
		Query            DiscoverQuery
		ExpectedResponse DiscoverResponse
	}{
		{
			Name: "Valid",
			Query: DiscoverQuery{
				EUI: basicstation.EUI{
					Prefix: "router",
					EUI64:  eui,
				},
			},
			ExpectedResponse: DiscoverResponse{
				EUI: basicstation.EUI{Prefix: "router", EUI64: eui},
				Muxs: basicstation.EUI{
					Prefix: "muxs",
				},
				URI: "wss://thethings.example.com:8887/traffic/eui-1111111111111111",
			},
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			msg, err := json.Marshal(tc.Query)
			a.So(err, should.BeNil)
			resp := th.HandleConnectionInfo(ctx, msg, mockServer, info, time.Now())
			expected, _ := json.Marshal(tc.ExpectedResponse)
			a.So(string(resp), should.Equal, string(expected))
		})
	}
}
