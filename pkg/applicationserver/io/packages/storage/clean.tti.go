// Copyright © 2020 The Things Industries B.V.

package storage

import (
	"github.com/mohae/deepcopy"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

// CleanApplicationUp removes ApplicationUp fields that we do not want to store.
func CleanApplicationUp(raw *ttnpb.ApplicationUp) *ttnpb.ApplicationUp {
	up := deepcopy.Copy(raw).(*ttnpb.ApplicationUp)
	up.JoinEUI = nil
	up.CorrelationIDs = nil

	switch up.Up.(type) {
	case *ttnpb.ApplicationUp_UplinkMessage:
		uplink, ok := up.Up.(*ttnpb.ApplicationUp_UplinkMessage)
		if !ok {
			return up
		}
		uplink.UplinkMessage.SessionKeyID = nil
		for _, md := range uplink.UplinkMessage.RxMetadata {
			md.UplinkToken = nil
		}
	}

	return up
}
