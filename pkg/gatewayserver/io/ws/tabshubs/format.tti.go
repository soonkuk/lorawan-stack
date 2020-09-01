// Copyright © 2020 The Things Industries B.V.

package tabshubs

import (
	"fmt"

	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/gatewayserver/io"
	"go.thethings.network/lorawan-stack/v3/pkg/gatewayserver/io/ws"
)

var (
	errSessionStateNotFound = errors.DefineUnavailable("session_state_not_found", "session state not found")
	trafficEndPointPrefix   = "/traffic"
)

// State represents the Tabs Hubs Session state.
type State struct {
	ID int32
}

type tabsHubs struct {
	tokens io.DownlinkTokens
}

// NewFormatter returns a new Tabs Hubs formatter.
func NewFormatter() ws.Formatter {
	return &tabsHubs{}
}

func (f *tabsHubs) Endpoints() ws.Endpoints {
	return ws.Endpoints{
		ConnectionInfo: "/router-info",
		Traffic:        fmt.Sprintf("%s/:id", trafficEndPointPrefix),
	}
}
