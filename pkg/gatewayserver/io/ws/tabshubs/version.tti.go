// Copyright © 2020 The Things Industries B.V.

package tabshubs

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.thethings.network/lorawan-stack/v3/pkg/frequencyplans"
	"go.thethings.network/lorawan-stack/v3/pkg/log"
	pfconfig "go.thethings.network/lorawan-stack/v3/pkg/pfconfig/tabshubs"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

// Version contains version information.
// This message is sent by the gateway.
type Version struct {
	Station  string `json:"station"`
	Firmware string `json:"firmware"`
}

// MarshalJSON implements json.Marshaler.
func (v Version) MarshalJSON() ([]byte, error) {
	type Alias Version
	return json.Marshal(struct {
		Type string `json:"msgtype"`
		Alias
	}{
		Type:  TypeUpstreamVersion,
		Alias: Alias(v),
	})
}

// GetRouterConfig gets router config for the particular version message.
func (f *tabsHubs) GetRouterConfig(ctx context.Context, msg []byte, bandID string, fp *frequencyplans.FrequencyPlan, receivedAt time.Time) (context.Context, []byte, *ttnpb.GatewayStatus, error) {
	var version Version
	if err := json.Unmarshal(msg, &version); err != nil {
		return nil, nil, nil, err
	}
	cfg, err := pfconfig.GetRouterConfig(bandID, fp, time.Now())
	if err != nil {
		return nil, nil, nil, err
	}
	routerCfg, err := cfg.MarshalJSON()
	if err != nil {
		return nil, nil, nil, err
	}
	// TODO: Revisit these fields for v3 events (https://github.com/TheThingsNetwork/lorawan-stack/issues/2629)
	stat := &ttnpb.GatewayStatus{
		Time: receivedAt,
		Versions: map[string]string{
			"station":  version.Station,
			"firmware": version.Firmware,
			"platform": fmt.Sprintf("TabsHubs - Firmware %s", version.Firmware),
		},
	}

	return log.NewContextWithFields(ctx, log.Fields(
		"station", version.Station,
		"firmware", version.Firmware,
	)), routerCfg, stat, nil
}
