// Copyright © 2020 The Things Industries B.V.

package storage

import "go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages"

// Option applies an option to the storage package
type Option func(p *storage)

// WithEnabled is an options that configures the enabled message types for the storage integration.
func WithEnabled(enabled packages.StorageTypesConfig) Option {
	return func(p *storage) {
		p.enabled = enabled
	}
}

var defaultOptions = []Option{}
