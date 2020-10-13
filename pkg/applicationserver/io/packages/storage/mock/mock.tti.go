// Copyright © 2020 The Things Industries B.V.

package mock

import (
	"context"

	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages/storage/provider"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

// Provider is a mock for the provider.Provider interface
type Provider struct {
	StoreCh chan []*io.ContextualApplicationUp
	QueryCh chan provider.Query
}

var (
	errFailForward = errors.DefineInternal("fail_forward", "failed to forward")
)

// New creates a new mock provider. Returns the provider, as well as the channels for tracking the incoming ApplicationUps, and GetApplicationUp requests.
func New(storeChSize int, queryChSize int) (*Provider, func()) {
	p := &Provider{
		StoreCh: make(chan []*io.ContextualApplicationUp, storeChSize),
		QueryCh: make(chan provider.Query, queryChSize),
	}
	closeFunc := func() {
		close(p.StoreCh)
		close(p.QueryCh)
	}

	return p, closeFunc
}

// Name implements the provider.Provider interface.
func (p *Provider) Name() string { return "mock" }

// Store implements the provider.Provider interface.
func (p *Provider) Store(ups []*io.ContextualApplicationUp) error {
	select {
	case p.StoreCh <- ups:
		return nil
	default:
		return errFailForward.New()
	}
}

// Range implements the provider.Provider interface.
func (p *Provider) Range(ctx context.Context, q provider.Query, f func(*ttnpb.ApplicationUp) error) error {
	select {
	case p.QueryCh <- q:
		return nil
	default:
		return errFailForward.New()
	}
}

// Close implements the provider.Provider interface.
func (p *Provider) Close() {}
