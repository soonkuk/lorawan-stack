// Copyright © 2020 The Things Industries B.V.

package postgres

import (
	"context"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Postgres database driver.
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages/storage/provider"
)

// Initialize sets up the database schemas.
func Initialize(db *gorm.DB) error {
	if err := db.AutoMigrate(models).Error; err != nil {
		return err
	}

	// Create hypertable for TimescaleDB
	var result struct{ Extname string }
	if err := db.Raw("SELECT extname FROM pg_extension WHERE extname = 'timescaledb'").Scan(&result).Error; err != nil && result.Extname == "timescaledb" {
		if err := db.Exec(fmt.Sprintf("SELECT create_hypertable('%s', 'timestamp', if_not_exists => true)", tableName)).Error; err != nil {
			return err
		}
	}
	return nil
}

// Open opens a new database connection based on config.
func Open(dsn string, config packages.PostgresConfig) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if config.Debug {
		db = db.Debug()
	}
	return db, nil
}

// New creates a new PostgreSQL provider for the storage integration.
func New(ctx context.Context, config packages.PostgresConfig) (provider.Provider, error) {
	p := &postgres{
		insertBatchSize: config.InsertBatchSize,
		selectBatchSize: config.SelectBatchSize,
	}

	if p.insertBatchSize < 0 {
		return nil, errNegativeBatchSize.WithAttributes("size", p.insertBatchSize)
	} else if p.insertBatchSize == 0 {
		p.insertBatchSize = 1024
	}
	if p.selectBatchSize < 0 {
		return nil, errNegativeBatchSize.WithAttributes("size", p.selectBatchSize)
	} else if p.selectBatchSize == 0 {
		p.selectBatchSize = 1024
	}

	var err error
	p.db, err = Open(config.DatabaseURI, config)
	if err != nil {
		return nil, err
	}
	go func() {
		select {
		case <-ctx.Done():
			p.db.Close()
		}
	}()
	if config.ReadDatabaseURI != "" {
		p.dbRead, err = Open(config.ReadDatabaseURI, config)
		if err != nil {
			return nil, err
		}
		go func() {
			select {
			case <-ctx.Done():
				p.dbRead.Close()
			}
		}()
	} else {
		p.dbRead = p.db
	}

	return p, nil
}
