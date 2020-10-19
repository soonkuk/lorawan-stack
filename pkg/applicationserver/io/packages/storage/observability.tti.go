// Copyright © 2020 The Things Industries B.V.

package storage

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"go.thethings.network/lorawan-stack/v3/pkg/events"
	"go.thethings.network/lorawan-stack/v3/pkg/metrics"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

var (
	evtStoreUp = events.Define(
		"as.packages.storage.up.store", "store upstream data message",
		events.WithVisibility(ttnpb.RIGHT_APPLICATION_TRAFFIC_READ),
	)
	evtDropUp = events.Define(
		"as.packages.storage.up.drop", "drop upstream data message",
		events.WithVisibility(ttnpb.RIGHT_APPLICATION_TRAFFIC_READ),
	)
)

const (
	subsystem     = "as_packages_storage"
	unknown       = "unknown"
	providerLabel = "provider"
)

var storageMetrics = &packageMetrics{
	upstreamMessagesStored: metrics.NewContextualCounterVec(
		prometheus.CounterOpts{
			Subsystem: subsystem,
			Name:      "upstream_messages_stored",
			Help:      "Number of upstream messages stored",
		},
		[]string{providerLabel},
	),
	upstreamMessagesDropped: metrics.NewContextualCounterVec(
		prometheus.CounterOpts{
			Subsystem: subsystem,
			Name:      "upstream_messages_dropped",
			Help:      "Number of upstream messages dropped",
		},
		[]string{providerLabel},
	),
}

func init() {
	metrics.MustRegister(storageMetrics)
}

type packageMetrics struct {
	upstreamMessagesStored  *metrics.ContextualCounterVec
	upstreamMessagesDropped *metrics.ContextualCounterVec
}

func (m packageMetrics) Describe(ch chan<- *prometheus.Desc) {
	m.upstreamMessagesStored.Describe(ch)
	m.upstreamMessagesDropped.Describe(ch)
}

func (m packageMetrics) Collect(ch chan<- prometheus.Metric) {
	m.upstreamMessagesStored.Collect(ch)
	m.upstreamMessagesDropped.Collect(ch)
}

func registerStoreUp(ctx context.Context, provider string, ids ttnpb.EndDeviceIdentifiers) {
	storageMetrics.upstreamMessagesStored.WithLabelValues(ctx, provider).Inc()
	events.Publish(evtStoreUp.New(ctx, events.WithIdentifiers(ids)))
}

func registerDropUp(ctx context.Context, provider string, ids ttnpb.EndDeviceIdentifiers, err error) {
	storageMetrics.upstreamMessagesDropped.WithLabelValues(ctx, provider).Inc()
	events.Publish(evtDropUp.NewWithIdentifiersAndData(ctx, ids, err))
}
