// Copyright © 2020 The Things Industries B.V.

package postgres

import (
	"context"
	"database/sql"

	"github.com/jinzhu/gorm"

	gormbulk "github.com/t-tiger/gorm-bulk-insert/v2"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages/storage/provider"
	"go.thethings.network/lorawan-stack/v3/pkg/jsonpb"
	"go.thethings.network/lorawan-stack/v3/pkg/log"
	"go.thethings.network/lorawan-stack/v3/pkg/tenant"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

const (
	providerName = "postgres"
)

// postgres implements the provider.Provider interface.
type postgres struct {
	db     *gorm.DB
	dbRead *gorm.DB

	insertBatchSize,
	selectBatchSize int
}

// Name implements the provider.Provider interface.
func (p *postgres) Name() string { return providerName }

// Store implements the provider.Provider interface.
func (p *postgres) Store(ups []*io.ContextualApplicationUp) error {
	insert := make([]interface{}, 0, len(ups))
	for _, up := range ups {
		b, err := jsonpb.TTN().Marshal(up.ApplicationUp)
		if err != nil {
			return err
		}
		toDB := StoredApplicationUp{
			Timestamp:     *up.ApplicationUp.ReceivedAt,
			Type:          applicationUpType(up.ApplicationUp),
			TenantID:      tenant.FromContext(up.Context).TenantID,
			DeviceID:      up.DeviceID,
			ApplicationID: up.ApplicationID,
			Data:          b,
		}
		if uplink := up.ApplicationUp.GetUplinkMessage(); uplink != nil {
			toDB.FPort = int(uplink.FPort)
		}
		insert = append(insert, toDB)
	}
	if err := gormbulk.BulkInsert(p.db, insert, p.insertBatchSize); err != nil {
		return errDatabase.WithCause(err)
	}
	return nil
}

// Range implements the provider.Provider interface.
func (p *postgres) Range(ctx context.Context, q provider.Query, f func(*ttnpb.ApplicationUp) error) error {
	db := p.dbRead

	db, err := withTenant(ctx, db)
	if err != nil {
		return err
	}
	db, err = withQuery(db, q)
	if err != nil {
		return err
	}

	db = db.BeginTx(ctx, &sql.TxOptions{
		ReadOnly:  true,
		Isolation: sql.LevelRepeatableRead,
	})
	defer func() {
		if err := db.Commit().Error; err != nil {
			log.FromContext(ctx).WithError(err).Error("Failed read-only transaction")
		}
	}()

	limit := q.Limit
	offset := 0
	storedUps := []StoredApplicationUp{}
	for {
		currentLimit := p.selectBatchSize
		if limit != nil && offset+currentLimit > int(limit.Value) {
			currentLimit = int(limit.Value) - offset
		}
		if currentLimit == 0 {
			return nil
		}
		if err := db.Offset(offset).Limit(currentLimit).Find(&storedUps).Error; err != nil {
			return errDatabase.WithCause(err)
		}
		for _, storedUp := range storedUps {
			up := &ttnpb.ApplicationUp{}
			if err := jsonpb.TTN().Unmarshal(storedUp.Data, up); err != nil {
				// TODO: Handle database corruption (https://github.com/TheThingsIndustries/lorawan-stack/issues/2383)
				continue
			}
			if err := f(up); err != nil {
				return err
			}
		}
		if len(storedUps) < currentLimit {
			return nil
		}
		offset += currentLimit
	}
}
