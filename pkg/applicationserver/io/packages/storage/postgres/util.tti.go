// Copyright © 2020 The Things Industries B.V.

package postgres

import (
	"time"

	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

func applicationUpType(up *ttnpb.ApplicationUp) string {
	if up == nil {
		return ""
	}

	switch up.Up.(type) {
	case *ttnpb.ApplicationUp_UplinkMessage:
		return "uplink_message"
	case *ttnpb.ApplicationUp_JoinAccept:
		return "join_accept"
	case *ttnpb.ApplicationUp_DownlinkAck:
		return "downlink_ack"
	case *ttnpb.ApplicationUp_DownlinkNack:
		return "downlink_nack"
	case *ttnpb.ApplicationUp_DownlinkSent:
		return "downlink_sent"
	case *ttnpb.ApplicationUp_DownlinkFailed:
		return "downlink_failed"
	case *ttnpb.ApplicationUp_DownlinkQueued:
		return "downlink_queued"
	case *ttnpb.ApplicationUp_DownlinkQueueInvalidated:
		return "downlink_queue_invalidated"
	case *ttnpb.ApplicationUp_LocationSolved:
		return "location_solved"
	case *ttnpb.ApplicationUp_ServiceData:
		return "service_data"
	default:
		return ""
	}
}

func nowPtr() *time.Time { t := time.Now().UTC(); return &t }
