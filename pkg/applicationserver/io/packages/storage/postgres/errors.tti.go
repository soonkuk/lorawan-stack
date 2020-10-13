// Copyright © 2020 The Things Industries B.V.

package postgres

import (
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
)

var (
	errDatabase   = errors.DefineInvalidArgument("database", "database error")
	errNoIDs      = errors.DefineInvalidArgument("no_device_id", "no device or application id set")
	errNoTenantID = errors.DefineInvalidArgument("no_tenant_id", "no tenant id set")

	errNegativeBatchSize = errors.DefineInvalidArgument("negative_batch_size", "batch size `{size}` is negative")
)
