// Copyright © 2020 The Things Industries B.V.

package storage_test

import (
	"testing"

	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages/storage"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
)

func TestClean(t *testing.T) {
	a := assertions.New(t)

	raw := ttnpb.NewPopulatedApplicationUp(test.Randy, false)
	clean := storage.CleanApplicationUp(raw)

	a.So(raw.JoinEUI, should.NotBeNil)
	a.So(raw.CorrelationIDs, should.NotBeNil)
	a.So(clean.JoinEUI, should.BeNil)
	a.So(clean.CorrelationIDs, should.BeNil)

	uplink := ttnpb.NewPopulatedApplicationUp_UplinkMessage(test.Randy, false)
	raw.Up = uplink

	clean = storage.CleanApplicationUp(raw)
	a.So(raw.JoinEUI, should.NotBeNil)
	a.So(raw.CorrelationIDs, should.NotBeNil)
	a.So(clean.JoinEUI, should.BeNil)
	a.So(clean.CorrelationIDs, should.BeNil)

	msg := clean.GetUplinkMessage()
	a.So(msg, should.NotBeNil)
	a.So(msg.SessionKeyID, should.BeNil)
	for _, md := range msg.RxMetadata {
		a.So(md.GetUplinkToken(), should.BeNil)
	}
}
