// Copyright © 2020 The Things Industries B.V.

package bulk

import (
	"context"
	"time"

	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages/storage/provider"
	"go.thethings.network/lorawan-stack/v3/pkg/log"
)

// New wraps a storage provider with bulk insert capabilities.
func New(ctx context.Context, p provider.Provider, q Queue, ticker <-chan time.Time) (provider.Provider, error) {
	bulk := &bulkProvider{
		ctx:    ctx,
		p:      p,
		q:      q,
		ticker: ticker,
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker:
				if err := bulk.flush(); err != nil {
					log.FromContext(bulk.ctx).WithError(err).Error("Failed to flush upstream messages")
				}
			}
		}
	}()
	return bulk, nil
}
