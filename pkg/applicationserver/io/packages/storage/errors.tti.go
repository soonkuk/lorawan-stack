// Copyright © 2020 The Things Industries B.V.

package storage

import "go.thethings.network/lorawan-stack/v3/pkg/errors"

var (
	errNotImplemented  = errors.DefineUnimplemented("not_implemented", "provider `{provider}` not implemented")
	errNoProvider      = errors.DefineUnavailable("no_provider", "no provider configured for storage")
	errInvalidInterval = errors.DefineInvalidArgument("invalid_interval", "invalid interval `{interval}`")

	errNoAppID = errors.DefineInvalidArgument("no_app_id", "no app id")
)
