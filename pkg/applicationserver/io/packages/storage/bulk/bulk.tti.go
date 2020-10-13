// Copyright © 2020 The Things Industries B.V.

package bulk

import (
	"context"
	"time"

	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages/storage/provider"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

const (
	providerName = "bulk"
)

// bulkProvider implements the provider.Provider interface.
type bulkProvider struct {
	ctx    context.Context
	p      provider.Provider
	q      Queue
	ticker <-chan time.Time
}

// Name implements the provider.Provider interface.
func (p *bulkProvider) Name() string {
	return p.p.Name()
}

// Store implements the provider.Provider interface.
func (p *bulkProvider) Store(ups []*io.ContextualApplicationUp) error {
	remaining := ups
	for {
		var err error
		if remaining, err = p.q.Push(remaining); err != nil {
			return err
		}
		if remaining == nil || len(remaining) == 0 {
			return nil
		}
		if err := p.flush(); err != nil {
			return err
		}
	}
}

// Range implements the provider.Provider interface.
func (p *bulkProvider) Range(ctx context.Context, q provider.Query, f func(*ttnpb.ApplicationUp) error) error {
	return p.p.Range(ctx, q, f)
}

func (p *bulkProvider) flush() error {
	ups, err := p.q.Pop(-1)
	if err != nil {
		return err
	}
	if len(ups) == 0 {
		return nil
	}
	return p.p.Store(ups)
}
