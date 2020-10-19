// Copyright © 2020 The Things Industries B.V.

package bulk_test

import (
	"testing"
	"time"

	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages/storage/bulk"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages/storage/mock"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages/storage/provider"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
)

func TestBulkProvider(t *testing.T) {
	ctx := test.Context()
	r := test.Randy
	a := assertions.New(t)
	q, err := bulk.NewMemoryQueue(2)
	a.So(err, should.BeNil)

	tickerCh := make(chan time.Time)
	p, close := mock.New(2, 1)
	defer close()
	b, err := bulk.New(ctx, p, q, tickerCh)
	a.So(err, should.BeNil)
	a.So(b, should.NotBeNil)

	ups := []*io.ContextualApplicationUp{newUp(r), newUp(r), newUp(r)}

	t.Run("Store", func(t *testing.T) {
		t.Run("Bulk", func(t *testing.T) {
			a := assertions.New(t)
			a.So(b.Store(ups[0:1]), should.BeNil)
			a.So(b.Store(ups[1:2]), should.BeNil)
			tickerCh <- time.Now()
			select {
			case result, ok := <-p.StoreCh:
				a.So(ok, should.BeTrue)
				a.So(result, should.Resemble, ups[:2])
			case <-time.After(time.Second):
				t.Fatal("Time out waiting for upstream messages")
				t.FailNow()
			}
		})

		t.Run("One", func(t *testing.T) {
			a := assertions.New(t)
			a.So(b.Store(ups[2:3]), should.BeNil)
			tickerCh <- time.Now()
			select {
			case result, ok := <-p.StoreCh:
				a.So(ok, should.BeTrue)
				a.So(result, should.Resemble, ups[2:3])
			case <-time.After(time.Second):
				t.Fatal("Time out waiting for upstream messages")
				t.FailNow()
			}
		})

		t.Run("Overflow", func(t *testing.T) {
			a := assertions.New(t)
			q.Pop(-1)
			a.So(b.Store(ups), should.BeNil)
			select {
			case result, ok := <-p.StoreCh:
				a.So(ok, should.BeTrue)
				a.So(result, should.Resemble, ups[0:2])
			case <-time.After(time.Second):
				t.Fatal("Time out waiting for upstream messages")
				t.FailNow()
			}

			tickerCh <- time.Now()

			select {
			case result, ok := <-p.StoreCh:
				a.So(ok, should.BeTrue)
				a.So(result, should.Resemble, ups[2:3])
			case <-time.After(time.Second):
				t.Fatal("Time out waiting for upstream messages")
				t.FailNow()
			}
		})
	})

	nullF := func(_ *ttnpb.ApplicationUp) error { return nil }
	t.Run("Range", func(t *testing.T) {
		req := ttnpb.NewPopulatedGetStoredApplicationUpRequest(r, false)
		q := provider.QueryFromRequest(req)
		a.So(p.Range(ctx, q, nullF), should.BeNil)

		select {
		case result, ok := <-p.QueryCh:
			a.So(ok, should.BeTrue)
			a.So(result, should.Resemble, q)
		case <-time.After(time.Second):
			t.Fatal("Time out waiting for upstream messages")
			t.FailNow()
		}
	})
}
