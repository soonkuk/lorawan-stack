// Copyright © 2020 The Things Industries B.V.

package postgres

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages/storage/provider"
	"go.thethings.network/lorawan-stack/v3/pkg/license"
	"go.thethings.network/lorawan-stack/v3/pkg/tenant"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

const tableName = "stored_application_ups"

var models = &StoredApplicationUp{}

// StoredApplicationUp is an upstream message stored in the database.
type StoredApplicationUp struct {
	Timestamp time.Time `gorm:"index;not null"`

	Type          string `gorm:"index;type:VARCHAR(36);not null"`
	DeviceID      string `gorm:"index;type:VARCHAR(36);not null"`
	ApplicationID string `gorm:"index;type:VARCHAR(36);not null"`
	TenantID      string `gorm:"index;type:VARCHAR(36);not null"`

	Data  []byte `gorm:"type:JSONB"`
	FPort int    `gorm:"index;type:INT;default:0"`
}

// TableName sets the table name
func (StoredApplicationUp) TableName() string {
	return tableName
}

func withDevID(db *gorm.DB, devID string) *gorm.DB {
	return db.Where("device_id = ?", devID)
}

func withAppID(db *gorm.DB, appID string) *gorm.DB {
	return db.Where("application_id = ?", appID)
}

func withQuery(readDB *gorm.DB, q provider.Query) (*gorm.DB, error) {
	db := readDB

	// Identifiers
	if ids := q.EndDeviceIDs; ids != nil {
		db = withDevID(db, ids.GetDeviceID())
		db = withAppID(db, ids.GetApplicationID())
	} else if ids := q.ApplicationIDs; ids != nil {
		db = withAppID(db, ids.GetApplicationID())
	} else {
		return nil, errNoIDs.New()
	}

	// Query
	switch q.Type {
	case "":
	default:
		if _, ok := ttnpb.StoredApplicationUpTypes[q.Type]; ok {
			db = db.Where("type = ?", q.Type)
		}
	}
	if q.Before != nil && q.After != nil {
		db = db.Where("timestamp BETWEEN ? AND ?", q.After, q.Before)
	} else if q.Before != nil {
		db = db.Where("timestamp < ?", q.Before)
	} else if q.After != nil {
		db = db.Where("timestamp > ?", q.After)
	}
	if fport := q.FPort; fport != nil {
		db = db.Where("f_port = ?", fport.Value)
	}

	// Order
	switch q.Order {
	case "received_at":
		db = db.Order("timestamp ASC")
	case "-received_at":
		db = db.Order("timestamp DESC")
	}

	return db, nil
}

func withTenant(ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	if !license.FromContext(ctx).MultiTenancy {
		return db, nil
	}

	tenantID := tenant.FromContext(ctx).TenantID
	if tenantID == "" {
		return nil, errNoTenantID
	}
	return db.Where("tenant_id = ?", tenantID), nil
}
