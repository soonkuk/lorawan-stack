// Copyright © 2020 The Things Industries B.V.

package provider

import (
	"context"
	"time"

	pbtypes "github.com/gogo/protobuf/types"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

// Provider stores and retrieves application layer messages.
type Provider interface {
	Name() string
	// Store stores upstream messages.
	Store(up []*io.ContextualApplicationUp) error
	// Range ranges over upstream messages and executes a callback function for each one.
	Range(ctx context.Context, q Query, f func(*ttnpb.ApplicationUp) error) error
}

// Query describes a query to the storage provider
type Query struct {
	ApplicationIDs *ttnpb.ApplicationIdentifiers
	EndDeviceIDs   *ttnpb.EndDeviceIdentifiers
	Limit, FPort   *pbtypes.UInt32Value
	Order, Type    string
	After, Before  *time.Time
}

// QueryFromRequest creates a new Query to the storage provider from a gRPC request.
func QueryFromRequest(req *ttnpb.GetStoredApplicationUpRequest) Query {
	return Query{
		ApplicationIDs: req.GetApplicationIDs(),
		EndDeviceIDs:   req.GetEndDeviceIDs(),
		Limit:          req.Limit,
		FPort:          req.FPort,
		Order:          req.Order,
		Type:           req.Type,
		After:          req.After,
		Before:         req.Before,
	}
}
