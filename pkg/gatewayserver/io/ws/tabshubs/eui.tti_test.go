// Copyright © 2020 The Things Industries B.V.

package tabshubs

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/types"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
)

func TestMarshalEUI(t *testing.T) {
	a := assertions.New(t)

	{
		eui := EUI{
			EUI64: types.EUI64{0xaa, 0xbb, 0x00, 0x01, 0x02, 0x03, 0x42, 0xff},
		}
		data, err := json.Marshal(eui)
		a.So(err, should.BeNil)
		a.So(string(data), should.Equal, `"aabb:1:203:42ff"`)
	}

	{
		eui := EUI{
			Prefix: "ROUTER",
			EUI64:  types.EUI64{0xaa, 0xbb, 0x00, 0x01, 0x02, 0x03, 0x42, 0xff},
		}
		data, err := json.Marshal(eui)
		a.So(err, should.BeNil)
		a.So(string(data), should.Equal, `"router-aabb:1:203:42ff"`)
	}

	{
		eui := EUI{
			Prefix: "muxs",
			EUI64:  types.EUI64{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		}
		data, err := json.Marshal(eui)
		a.So(err, should.BeNil)
		a.So(string(data), should.Equal, `"muxs-::0"`)
	}
}

func TestUnmarshalEUI(t *testing.T) {
	for i, tc := range []struct {
		Input  string
		Prefix string
		EUI64  types.EUI64
		OK     bool
	}{
		{
			Input: `"aa-bb-cc-01-02-03-42-ff"`,
			EUI64: types.EUI64{0xaa, 0xbb, 0xcc, 0x01, 0x02, 0x03, 0x42, 0xff},
			OK:    true,
		},
		{
			Input: `"aa:bb:cc:01:02:03:42:ff"`,
			EUI64: types.EUI64{0xaa, 0xbb, 0xcc, 0x01, 0x02, 0x03, 0x42, 0xff},
			OK:    true,
		},
		{
			Input: `"aa:bb:cc:01:02:03"`,
			OK:    false,
		},
		{
			Input: `aa:bb:cc:01:02:03:42:ff:f2`,
			OK:    false,
		},
		{
			Input: `aa:bb:cc:01:02:03:42:xx`,
			OK:    false,
		},
		{
			Input: `aa:bb:cc:01:02:03:42-01`,
			OK:    false,
		},
		{
			Input: `"aa-bb-cc-01-02-03"`,
			OK:    false, // Too short.
		},
		{
			Input: `aa-bb-cc-01-02-03-42-ff`,
			OK:    false, // Not a string.
		},
		{
			Input: `"aa-bb-cc-01-02-03-42-xx"`,
			OK:    false, // Invalid hex.
		},
		{
			Input: `"aabb:cc01:0203:42ff"`,
			EUI64: types.EUI64{0xaa, 0xbb, 0xcc, 0x01, 0x02, 0x03, 0x42, 0xff},
			OK:    true,
		},
		{
			Input: `"aabb:01:203:42ff"`,
			EUI64: types.EUI64{0xaa, 0xbb, 0x00, 0x01, 0x02, 0x03, 0x42, 0xff},
			OK:    true,
		},
		{
			Input: `"aabb:01::"`,
			EUI64: types.EUI64{0xaa, 0xbb, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00},
			OK:    true,
		},
		{
			Input:  `"router-aabb:01::"`,
			Prefix: "router",
			EUI64:  types.EUI64{0xaa, 0xbb, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00},
			OK:     true,
		},
		{
			Input: `"::0"`,
			EUI64: types.EUI64{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			OK:    true,
		},
		{
			Input:  `"muxs-::0"`,
			Prefix: "muxs",
			EUI64:  types.EUI64{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			OK:     true,
		},
		{
			Input: `12302426811387609088`,
			EUI64: types.EUI64{0xaa, 0xbb, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00},
			OK:    true,
		},
		{
			Input: `-12302426811387609088`,
			OK:    false,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			a := assertions.New(t)
			var eui EUI
			err := json.Unmarshal([]byte(tc.Input), &eui)
			if tc.OK {
				a.So(err, should.BeNil)
				a.So(eui.EUI64, should.Resemble, tc.EUI64)
			} else {
				a.So(err, should.NotBeNil)
			}
		})
	}
}
