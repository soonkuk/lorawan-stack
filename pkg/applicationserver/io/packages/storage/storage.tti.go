// Copyright © 2020 The Things Industries B.V.

package storage

import (
	"context"

	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

// HandleUp implements packages.ApplicationPackageHandler.
func (p *storage) HandleUp(ctx context.Context, def *ttnpb.ApplicationPackageDefaultAssociation, assoc *ttnpb.ApplicationPackageAssociation, raw *ttnpb.ApplicationUp) error {
	if def == nil && assoc == nil {
		return nil
	}

	if !p.shouldStore(raw) {
		return nil
	}
	up := CleanApplicationUp(raw)
	err := p.provider.Store([]*io.ContextualApplicationUp{{
		Context:       ctx,
		ApplicationUp: up,
	}})
	if err != nil {
		registerDropUp(ctx, p.provider.Name(), up.EndDeviceIdentifiers, err)
		return err
	}
	registerStoreUp(ctx, p.provider.Name(), up.EndDeviceIdentifiers)
	return nil
}

// Package implements packages.ApplicationPackageHandler.
func (p *storage) Package() *ttnpb.ApplicationPackage {
	return &ttnpb.ApplicationPackage{
		Name:         packageName,
		DefaultFPort: 200,
	}
}

// shouldStore checks if the storage package is configured for storing the message based on its type.
func (p *storage) shouldStore(up *ttnpb.ApplicationUp) bool {
	if p.enabled.All {
		return true
	}

	switch up.Up.(type) {
	case *ttnpb.ApplicationUp_UplinkMessage:
		return p.enabled.UplinkMessage
	case *ttnpb.ApplicationUp_JoinAccept:
		return p.enabled.JoinAccept
	case *ttnpb.ApplicationUp_DownlinkAck:
		return p.enabled.DownlinkAck
	case *ttnpb.ApplicationUp_DownlinkNack:
		return p.enabled.DownlinkNack
	case *ttnpb.ApplicationUp_DownlinkFailed:
		return p.enabled.DownlinkFailed
	case *ttnpb.ApplicationUp_DownlinkSent:
		return p.enabled.DownlinkSent
	case *ttnpb.ApplicationUp_DownlinkQueued:
		return p.enabled.DownlinkQueued
	case *ttnpb.ApplicationUp_DownlinkQueueInvalidated:
		return p.enabled.DownlinkQueueInvalidated
	case *ttnpb.ApplicationUp_LocationSolved:
		return p.enabled.LocationSolved
	case *ttnpb.ApplicationUp_ServiceData:
		return p.enabled.ServiceData
	}

	return false
}
