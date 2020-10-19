// Copyright © 2020 The Things Industries B.V.

package bulk_test

import (
	"testing"

	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages/storage/bulk"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/util/randutil"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
)

func newUp(r *randutil.LockedRand) *io.ContextualApplicationUp {
	return &io.ContextualApplicationUp{
		Context:       test.Context(),
		ApplicationUp: ttnpb.NewPopulatedApplicationUp(r, true),
	}
}

func TestMemoryQueue(t *testing.T) {
	r := test.Randy
	a := assertions.New(t)

	q, err := bulk.NewMemoryQueue(3)
	a.So(err, should.BeNil)
	a.So(q, should.NotBeNil)

	ups := []*io.ContextualApplicationUp{newUp(r), newUp(r), newUp(r), newUp(r)}

	t.Run("PopOne", func(t *testing.T) {
		a := assertions.New(t)
		remaining, err := q.Push(ups[0:1])
		a.So(remaining, should.BeNil)
		a.So(err, should.BeNil)

		res, err := q.Pop(1)
		a.So(err, should.BeNil)
		a.So(res, should.Resemble, ups[0:1])

		res, err = q.Pop(1)
		a.So(err, should.BeNil)
		a.So(len(res), should.Equal, 0)
	})

	t.Run("PopMore", func(t *testing.T) {
		a := assertions.New(t)
		remaining, err := q.Push(ups[0:1])
		a.So(remaining, should.BeNil)
		a.So(err, should.BeNil)

		res, err := q.Pop(3)
		a.So(err, should.BeNil)
		a.So(res, should.Resemble, ups[0:1])

		res, err = q.Pop(1)
		a.So(err, should.BeNil)
		a.So(len(res), should.Equal, 0)
	})

	t.Run("PopAll", func(t *testing.T) {
		a := assertions.New(t)
		remaining, err := q.Push(ups[0:2])
		a.So(remaining, should.BeNil)
		a.So(err, should.BeNil)

		res, err := q.Pop(-1)
		a.So(err, should.BeNil)
		a.So(res, should.Resemble, ups[0:2])

		res, err = q.Pop(1)
		a.So(err, should.BeNil)
		a.So(len(res), should.Equal, 0)
	})

	t.Run("PopSome", func(t *testing.T) {
		a := assertions.New(t)
		remaining, err := q.Push(ups[0:2])
		a.So(remaining, should.BeNil)
		a.So(err, should.BeNil)

		res, err := q.Pop(1)
		a.So(err, should.BeNil)
		a.So(res, should.Resemble, ups[0:1])

		res, err = q.Pop(1)
		a.So(err, should.BeNil)
		a.So(res, should.Resemble, ups[1:2])
	})

	t.Run("Overflow", func(t *testing.T) {
		a := assertions.New(t)

		remaining, err := q.Push(ups)
		a.So(remaining, should.Resemble, []*io.ContextualApplicationUp{ups[3]})
		a.So(err, should.BeNil)

		res, err := q.Pop(3)
		a.So(err, should.BeNil)
		a.So(res, should.Resemble, ups[0:3])
	})
}
