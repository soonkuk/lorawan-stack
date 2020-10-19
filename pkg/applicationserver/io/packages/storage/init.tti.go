// Copyright © 2020 The Things Industries B.V.

package storage

import (
	"context"
	"time"

	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages/storage/bulk"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages/storage/postgres"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages/storage/provider"
	"go.thethings.network/lorawan-stack/v3/pkg/log"
)

const packageName = "storage-integration"

// storage is the Storage Integration package.
type storage struct {
	ctx      context.Context
	provider provider.Provider
	enabled  packages.StorageTypesConfig
}

// FromConfig instantiates the Storage Integration package from configuration.
func FromConfig(ctx context.Context, config packages.StorageConfig) (packages.ApplicationPackageHandler, error) {
	sctx := log.NewContextWithField(ctx, "namespace", "applicationserver/io/packages/storage")

	p, err := initProvider(ctx, config)
	if err != nil {
		return nil, err
	}

	if config.Bulk.Enabled {
		if config.Bulk.Interval < 0 {
			return nil, errInvalidInterval.WithAttributes("interval", config.Bulk.Interval)
		} else if config.Bulk.Interval == 0 {
			config.Bulk.Interval = 10 * time.Second
		}
		queue, err := bulk.NewMemoryQueue(config.Bulk.MaxSize)
		if err != nil {
			return nil, err
		}
		p, err = bulk.New(sctx, p, queue, time.NewTicker(config.Bulk.Interval).C)
		if err != nil {
			return nil, err
		}
	}
	return New(sctx, p, WithEnabled(config.Enable))
}

// New instantiates the Storage Integration package.
func New(ctx context.Context, provider provider.Provider, options ...Option) (packages.ApplicationPackageHandler, error) {
	log.FromContext(ctx).WithFields(log.Fields(
		"provider", provider.Name(),
	)).Info("Initialized Storage Integration")

	p := &storage{
		ctx:      ctx,
		provider: provider,
	}
	opts := append(defaultOptions, options...)
	for _, o := range opts {
		o(p)
	}
	return p, nil
}

// initProvider initializes the storage provider.
func initProvider(ctx context.Context, cfg packages.StorageConfig) (provider.Provider, error) {
	switch cfg.Provider {
	case "":
		return nil, errNoProvider.New()
	case "postgres":
		return postgres.New(ctx, cfg.Postgres)
	default:
		return nil, errNotImplemented.WithAttributes("provider", cfg.Provider)
	}
}
