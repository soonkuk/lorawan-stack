// Copyright © 2020 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package packages

import "time"

// Config contains configuration options for application packages.
type Config struct {
	Storage StorageConfig `name:"storage"`
}

// StorageConfig contains configuration options for the storage integration.
type StorageConfig struct {
	Provider string             `name:"provider" description:"Storage provider (postgres)"`
	Enable   StorageTypesConfig `name:"enable"`
	Bulk     BulkConfig         `name:"bulk"`
	Postgres PostgresConfig     `name:"postgres"`
}

// StorageTypesConfig contains configuration options for which message types will be stored.
type StorageTypesConfig struct {
	All                      bool `name:"all" description:"Store all existing and future message types"`
	UplinkMessage            bool `name:"uplink-message" description:"Store uplink messages"`
	JoinAccept               bool `name:"join-accept" description:"Store join accept messages"`
	DownlinkAck              bool `name:"downlink-ack" description:"Store downlink ack messages"`
	DownlinkNack             bool `name:"downlink-nack" description:"Store downlink nack messages"`
	DownlinkSent             bool `name:"downlink-sent" description:"Store downlink sent messages"`
	DownlinkFailed           bool `name:"downlink-failed" description:"Store downlink failed messages"`
	DownlinkQueued           bool `name:"downlink-failed" description:"Store downlink queued messages"`
	DownlinkQueueInvalidated bool `name:"downlink-failed" description:"Store downlink queue invalidated messages"`
	LocationSolved           bool `name:"location-solved" description:"Store location solved messages"`
	ServiceData              bool `name:"service-data" description:"Store service data messages"`
}

// BulkConfig contains configuration options for bulk storing uplinks.
type BulkConfig struct {
	Enabled  bool          `name:"enabled" description:"Store uplinks in batches"`
	MaxSize  int           `name:"max-size" description:"Max number of uplinks to store in cache"`
	Interval time.Duration `name:"interval" description:"Interval between storing uplinks"`
}

// PostgresConfig contains configuration options for the PostgreSQL storage provider.
type PostgresConfig struct {
	Debug           bool   `name:"debug" description:"Start in debug mode"`
	DatabaseURI     string `name:"database-uri" description:"Database connection URI (when provider is postgres)"`
	ReadDatabaseURI string `name:"read-database-uri" description:"Read-Only database connection URI"`
	InsertBatchSize int    `name:"insert-batch-size" description:"Batch size for INSERT commands"`
	SelectBatchSize int    `name:"select-batch-size" description:"Batch size for SELECT commands"`
}
