// Copyright © 2020 The Things Industries B.V.

package tabshubs

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.thethings.network/lorawan-stack/v3/pkg/band"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/frequencyplans"
	"go.thethings.network/lorawan-stack/v3/pkg/pfconfig/shared"
)

// DataRates encodes the available datarates of the channel plan for the Tabs Hubs in the format below:
// [0] -> SF (Spreading Factor; Range: 7...12 for LoRa, 0 for FSK)
// [1] -> BW (Bandwidth; 125/250/500 for LoRa, ignored for FSK)
// [2] -> DNONLY (Downlink Only; 1 = true, 0 = false)
type DataRates [16][3]int

// Upchannels encodes the uplink channels
// [0] -> Frequency (MHz)
// [1] -> Min Data Rate (number)
// [2] -> Max Data Rate (number)
type Upchannels [8][3]int

// RouterConfig contains the router configuration.
// This message is sent by the Gateway Server.
type RouterConfig struct {
	NetID          []int                 `json:"NetID"`
	JoinEUI        [][]int               `json:"JoinEui"`
	Region         string                `json:"region"`
	HardwareSpec   string                `json:"hwspec"`
	FrequencyRange []int                 `json:"freq_range"`
	DataRates      DataRates             `json:"DRs"`
	SX1301Config   []shared.SX1301Config `json:"sx1301_conf"`
	UpChannels     [][]int               `json:"upchannels"`
	RegionID       int                   `json:"regionid"`
	MuxTime        float64               `json:"MuxTime"`
	MaxEIRP        int                   `json:"max_eirp"`
	Protocol       int                   `json:"protocol"`
	Config         struct {
		Region string `json:"region"`
	} `json:"config"`
	Beaconing []int `json:"bcning"`
}

const (
	configHardwareSpecPrefix = "sx1301"
	configProtocol           = 1
)

var errFrequencyPlan = errors.DefineInvalidArgument("frequency_plan", "invalid frequency plan `{name}`")

// RegionIDToRegion maps the regionid field to the region field.
var RegionIDToRegion = map[string]int{
	"EU863": 1,
	"US902": 2,
	"EU433": 3,
	"AU915": 4,
	"CN470": 5,
	"CN779": 6,
	"AS923": 7,
	"KR920": 8,
	"IN865": 9,
	"IL915": 10,
	"RU864": 11,
}

// MarshalJSON implements json.Marshaler.
func (conf RouterConfig) MarshalJSON() ([]byte, error) {
	type Alias RouterConfig
	return json.Marshal(struct {
		Type string `json:"msgtype"`
		Alias
	}{
		Type:  "router_config",
		Alias: Alias(conf),
	})
}

// GetRouterConfig returns the routerconfig message to be sent to the gateway.
// Does not support multiple frequency plans.
func GetRouterConfig(bandID string, fp *frequencyplans.FrequencyPlan, dlTime time.Time) (RouterConfig, error) {
	if err := fp.Validate(); err != nil {
		return RouterConfig{}, errFrequencyPlan.New()
	}
	conf := RouterConfig{}
	conf.JoinEUI = nil
	conf.NetID = nil
	conf.Beaconing = nil

	phy, err := band.GetByID(bandID)
	if err != nil {
		return RouterConfig{}, errFrequencyPlan.New()
	}
	s := strings.Split(phy.ID, "_")
	if len(s) < 2 {
		return RouterConfig{}, errFrequencyPlan.New()
	}
	conf.Region = fmt.Sprintf("%s%s", s[0], s[1])
	conf.Config.Region = fmt.Sprintf("%s%s/tracknet", s[0], s[1])

	if len(fp.Radios) == 0 {
		return RouterConfig{}, errFrequencyPlan.New()
	}
	conf.FrequencyRange = []int{int(fp.Radios[0].TxConfiguration.MinFrequency), int(fp.Radios[0].TxConfiguration.MaxFrequency)}

	conf.MaxEIRP = int(phy.DefaultMaxEIRP)

	conf.Protocol = configProtocol

	conf.HardwareSpec = fmt.Sprintf("%s/%d", configHardwareSpecPrefix, 1)

	conf.DataRates, err = getDataRatesFromBandID(bandID)
	if err != nil {
		return RouterConfig{}, errFrequencyPlan.New()
	}

	for _, channel := range fp.UplinkChannels {
		upChannel := make([]int, 3)
		upChannel[0] = int(channel.Frequency)
		upChannel[1] = int(channel.MinDataRate)
		upChannel[2] = int(channel.MaxDataRate)
		conf.UpChannels = append(conf.UpChannels, upChannel)
	}

	sx1301Conf, err := shared.BuildSX1301Config(fp)
	// These fields are not defined in the v1.5 ref design https://doc.sm.tc/station/gw_v1.5.html#rfconf-object and would cause a parsing error.
	sx1301Conf.Radios[0].TxFreqMin = 0
	sx1301Conf.Radios[0].TxFreqMax = 0
	// Remove hardware specific values that are not necessary.
	sx1301Conf.TxLUTConfigs = nil
	for i := range sx1301Conf.Radios {
		sx1301Conf.Radios[i].Type = ""
	}
	if err != nil {
		return RouterConfig{}, err
	}
	conf.SX1301Config = append(conf.SX1301Config, *sx1301Conf)

	// Add the MuxTime for RTT measurement.
	conf.MuxTime = float64(dlTime.Unix()) + float64(dlTime.Nanosecond())/(1e9)

	return conf, nil
}

// getDataRatesFromBandID parses the available data rates from the band into DataRates.
func getDataRatesFromBandID(id string) (DataRates, error) {
	phy, err := band.GetByID(id)
	if err != nil {
		return DataRates{}, err
	}

	// Set the default values.
	drs := DataRates{}
	for _, dr := range drs {
		dr[0] = -1
		dr[1] = 0
		dr[2] = 0
	}

	for i, dr := range phy.DataRates {
		if loraDR := dr.Rate.GetLoRa(); loraDR != nil {
			loraDR.GetSpreadingFactor()
			drs[i][0] = int(loraDR.GetSpreadingFactor())
			drs[i][1] = int(loraDR.GetBandwidth() / 1000)
		} else if fskDR := dr.Rate.GetFSK(); fskDR != nil {
			drs[i][0] = 0 // must be set to 0 for FSK, the BW field is ignored.
		}
	}
	return drs, nil
}
