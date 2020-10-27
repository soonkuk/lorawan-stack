// Copyright © 2020 The Things Industries B.V.

package tabshubs

import (
	"testing"
	"time"

	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/frequencyplans"
	"go.thethings.network/lorawan-stack/v3/pkg/pfconfig/shared"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
)

func TestGetRouterConfig(t *testing.T) {
	for _, tc := range []struct {
		Name            string
		FrequencyPlan   frequencyplans.FrequencyPlan
		FrequencyPlanID string
		Cfg             RouterConfig
		IsProd          bool
		ErrorAssertion  func(err error) bool
	}{
		{
			Name:            "NilFrequencyPlan",
			FrequencyPlan:   frequencyplans.FrequencyPlan{},
			FrequencyPlanID: "dummy",
			Cfg:             RouterConfig{},
			ErrorAssertion: func(err error) bool {
				return errors.Resemble(err, errFrequencyPlan)
			},
		},
		{
			Name: "InvalidFrequencyPlan",
			FrequencyPlan: frequencyplans.FrequencyPlan{
				BandID: "PinkFloyd",
			},
			FrequencyPlanID: "dummy",
			Cfg:             RouterConfig{},
			ErrorAssertion: func(err error) bool {
				return errors.Resemble(err, errFrequencyPlan)
			},
		},
		{
			Name: "ValidFrequencyPlan",
			FrequencyPlan: frequencyplans.FrequencyPlan{
				BandID: "US_902_928",
				Radios: []frequencyplans.Radio{
					{
						Enable:    true,
						ChipType:  "SX1257",
						Frequency: 922300000,
						TxConfiguration: &frequencyplans.RadioTxConfiguration{
							MinFrequency: 909000000,
							MaxFrequency: 925000000,
						},
					},
					{
						Enable:    false,
						ChipType:  "SX1257",
						Frequency: 923000000,
					},
				},
			},
			FrequencyPlanID: "US_902_928",
			Cfg: RouterConfig{
				Region:       "US902",
				HardwareSpec: "sx1301/1",
				RegionID:     RegionToRegionID["US902"],
				MaxEIRP:      30,
				Protocol:     1,
				Config: Config{
					Region: "US902/tracknet",
				},
				FrequencyRange: []int{909000000, 925000000},
				DataRates: DataRates{
					[3]int{10, 125, 0},
					[3]int{9, 125, 0},
					[3]int{8, 125, 0},
					[3]int{7, 125, 0},
					[3]int{8, 500, 0},
					[3]int{0, 0, 0},
					[3]int{0, 0, 0},
					[3]int{0, 0, 0},
					[3]int{12, 500, 0},
					[3]int{11, 500, 0},
					[3]int{10, 500, 0},
					[3]int{9, 500, 0},
					[3]int{8, 500, 0},
					[3]int{7, 500, 0},
				},
				SX1301Config: []shared.SX1301Config{
					{
						LoRaWANPublic: true,
						ClockSource:   0,
						AntennaGain:   0,
						Radios: []shared.RFConfig{
							{
								Enable:    true,
								Frequency: 922300000,
							},
							{
								Enable:    false,
								Frequency: 923000000,
							},
						},
						Channels: []shared.IFConfig{
							{Enable: false, Radio: 0, IFValue: 0, Bandwidth: 0, SpreadFactor: 0, Datarate: 0},
							{Enable: false, Radio: 0, IFValue: 0, Bandwidth: 0, SpreadFactor: 0, Datarate: 0},
							{Enable: false, Radio: 0, IFValue: 0, Bandwidth: 0, SpreadFactor: 0, Datarate: 0},
							{Enable: false, Radio: 0, IFValue: 0, Bandwidth: 0, SpreadFactor: 0, Datarate: 0},
							{Enable: false, Radio: 0, IFValue: 0, Bandwidth: 0, SpreadFactor: 0, Datarate: 0},
							{Enable: false, Radio: 0, IFValue: 0, Bandwidth: 0, SpreadFactor: 0, Datarate: 0},
							{Enable: false, Radio: 0, IFValue: 0, Bandwidth: 0, SpreadFactor: 0, Datarate: 0},
							{Enable: false, Radio: 0, IFValue: 0, Bandwidth: 0, SpreadFactor: 0, Datarate: 0},
						},
						LoRaStandardChannel: &shared.IFConfig{Enable: false, Radio: 0, IFValue: 0, Bandwidth: 0, SpreadFactor: 0, Datarate: 0},
						FSKChannel:          &shared.IFConfig{Enable: false, Radio: 0, IFValue: 0, Bandwidth: 0, SpreadFactor: 0, Datarate: 0},
					},
				},
				UpChannels: nil,
			},
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			a := assertions.New(t)
			cfg, err := GetRouterConfig(tc.FrequencyPlan.BandID, &tc.FrequencyPlan, time.Now())
			if err != nil {
				if tc.ErrorAssertion == nil || !a.So(tc.ErrorAssertion(err), should.BeTrue) {
					t.Fatalf("Unexpected error: %v", err)
				}
			} else if tc.ErrorAssertion != nil {
				t.Fatalf("Expected error")
			} else {
				cfg.MuxTime = 0
				if !a.So(cfg, should.Resemble, tc.Cfg) {
					t.Fatalf("Invalid config: %v", cfg)
				}
			}
		})
	}
}

func TestGetDataRatesFromFrequencyPlan(t *testing.T) {
	for _, tc := range []struct {
		Name           string
		BandID         string
		DataRates      DataRates
		ErrorAssertion func(error) bool
	}{
		{
			Name:           "InvalidBandID",
			BandID:         "EU",
			DataRates:      DataRates{},
			ErrorAssertion: errors.IsNotFound,
		},
		{
			Name:   "ValidBAndID",
			BandID: "EU_433",
			DataRates: DataRates{
				[3]int{12, 125, 0},
				[3]int{11, 125, 0},
				[3]int{10, 125, 0},
				[3]int{9, 125, 0},
				[3]int{8, 125, 0},
				[3]int{7, 125, 0},
				[3]int{7, 250, 0},
			},
		},
		{
			Name:   "ValidBAndIDUS",
			BandID: "US_902_928",
			DataRates: DataRates{
				[3]int{10, 125, 0},
				[3]int{9, 125, 0},
				[3]int{8, 125, 0},
				[3]int{7, 125, 0},
				[3]int{8, 500, 0},
				[3]int{0, 0, 0},
				[3]int{0, 0, 0},
				[3]int{0, 0, 0},
				[3]int{12, 500, 0},
				[3]int{11, 500, 0},
				[3]int{10, 500, 0},
				[3]int{9, 500, 0},
				[3]int{8, 500, 0},
				[3]int{7, 500, 0},
			},
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			a := assertions.New(t)
			drs, err := getDataRatesFromBandID(tc.BandID)
			if err != nil && (tc.ErrorAssertion == nil || !a.So(tc.ErrorAssertion(err), should.BeTrue)) {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !a.So(drs, should.Resemble, tc.DataRates) {
				t.Fatalf("Invalid datarates: %v", drs)
			}
		})
	}
}
