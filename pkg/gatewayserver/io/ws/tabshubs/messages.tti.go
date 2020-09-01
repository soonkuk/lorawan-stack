// Copyright © 2020 The Things Industries B.V.

package tabshubs

import (
	"encoding/json"

	"go.thethings.network/lorawan-stack/v3/pkg/basicstation"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
)

// MessageType is the type of the message.
type MessageType string

// Definition of the Tabs Hubs message types.
const (
	// Upstream types for messages from the Gateway.
	TypeUpstreamVersion         = "version"
	TypeUpstreamJoinRequest     = "jreq"
	TypeUpstreamUplinkDataFrame = "updf"
	TypeUpstreamTxConfirmation  = "dntxed"

	// Downstream types for messages from the Network
	TypeDownstreamDownlinkMessage = "dnframe"
)

// DiscoverQuery contains the unique identifier of the gateway.
// This message is sent by the gateway.
type DiscoverQuery struct {
	EUI basicstation.EUI `json:"router"`
}

// DiscoverResponse contains the response to the discover query.
// This message is sent by the Gateway Server.
type DiscoverResponse struct {
	EUI   basicstation.EUI `json:"router"`
	Muxs  basicstation.EUI `json:"muxs,omitempty"`
	URI   string           `json:"uri,omitempty"`
	Error string           `json:"error,omitempty"`
}

var errNotSupported = errors.DefineFailedPrecondition("not_supported", "not supported")

// Type returns the message type of the given data.
func Type(data []byte) (string, error) {
	msg := struct {
		Type string `json:"msgtype"`
	}{}
	if err := json.Unmarshal(data, &msg); err != nil {
		return "", err
	}
	return msg.Type, nil
}
