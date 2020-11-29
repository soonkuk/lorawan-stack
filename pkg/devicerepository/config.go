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

package devicerepository

import (
	"context"

	"go.thethings.network/lorawan-stack/v3/pkg/devicerepository/store"
)

// Config represents the DeviceRepository configuration.
type Config struct {
	Static    map[string][]byte `name:"-"`
	Directory string            `name:"directory" description:"Retrieve devices from the filesystem"`
	URL       string            `name:"url" description:"Retrieve devices from a web server"`

	PhotosBaseURL string `name:"photos-base-url" description:"The base URL for photos assets"`
}

// NewStore creates a new Store for end devices.
func (c Config) NewStore(ctx context.Context) (store.Store, error) {
	return &store.NoopStore{}, nil
}
