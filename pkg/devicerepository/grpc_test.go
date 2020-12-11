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

package devicerepository_test

import (
	"testing"

	"go.thethings.network/lorawan-stack/v3/pkg/component"
	componenttest "go.thethings.network/lorawan-stack/v3/pkg/component/test"
	"go.thethings.network/lorawan-stack/v3/pkg/config"
	"go.thethings.network/lorawan-stack/v3/pkg/devicerepository"
	"go.thethings.network/lorawan-stack/v3/pkg/devicerepository/store"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
)

type mockStore struct {
	// last requests
	lastGetBrandsRequest store.GetBrandsRequest
	lastGetModelsRequest store.GetModelsRequest
	lastVersionIDs       *ttnpb.EndDeviceVersionIdentifiers

	// mock responses
	brands   []*ttnpb.EndDeviceBrand
	models   []*ttnpb.EndDeviceModel
	template *ttnpb.EndDeviceTemplate
	uplinkDecoder,
	downlinkDecoder,
	downlinkEncoder string

	// mock errors
	err error
}

// GetBrands lists available end device vendors.
func (s *mockStore) GetBrands(req store.GetBrandsRequest) (*store.GetBrandsResponse, error) {
	s.lastGetBrandsRequest = req
	if s.err != nil {
		return nil, s.err
	}
	if s.brands == nil {
		s.brands = []*ttnpb.EndDeviceBrand{}
	}
	return &store.GetBrandsResponse{
		Count:  uint32(len(s.brands)),
		Offset: 0,
		Total:  uint32(len(s.brands)),
		Brands: s.brands,
	}, nil
}

// GetModels lists available end device definitions.
func (s *mockStore) GetModels(req store.GetModelsRequest) (*store.GetModelsResponse, error) {
	s.lastGetModelsRequest = req
	if s.err != nil {
		return nil, s.err
	}
	if s.models == nil {
		s.models = []*ttnpb.EndDeviceModel{}
	}
	return &store.GetModelsResponse{
		Count:  uint32(len(s.models)),
		Offset: 0,
		Total:  uint32(len(s.models)),
		Models: s.models,
	}, nil
}

// GetTemplate retrieves an end device template for an end device definition.
func (s *mockStore) GetTemplate(ids *ttnpb.EndDeviceVersionIdentifiers) (*ttnpb.EndDeviceTemplate, error) {
	s.lastVersionIDs = ids
	return s.template, s.err
}

// GetUplinkDecoder retrieves the codec for decoding uplink messages.
func (s *mockStore) GetUplinkDecoder(ids *ttnpb.EndDeviceVersionIdentifiers) (string, error) {
	s.lastVersionIDs = ids
	return s.uplinkDecoder, s.err
}

// GetDownlinkDecoder retrieves the codec for decoding downlink messages.
func (s *mockStore) GetDownlinkDecoder(ids *ttnpb.EndDeviceVersionIdentifiers) (string, error) {
	s.lastVersionIDs = ids
	return s.downlinkDecoder, s.err
}

// GetDownlinkEncoder retrieves the codec for encoding downlink messages.
func (s *mockStore) GetDownlinkEncoder(ids *ttnpb.EndDeviceVersionIdentifiers) (string, error) {
	s.lastVersionIDs = ids
	return s.downlinkEncoder, s.err
}

func TestGRPC(t *testing.T) {
	ids := &ttnpb.EndDeviceVersionIdentifiers{
		BrandID:         "brand",
		ModelID:         "model",
		FirmwareVersion: "1.0",
		HardwareVersion: "1.0",
		BandID:          "band",
	}

	componentConfig := &component.Config{
		ServiceBase: config.ServiceBase{
			GRPC: config.GRPC{
				Listen:                      ":0",
				AllowInsecureForCredentials: true,
			},
		},
	}
	c := componenttest.NewComponent(t, componentConfig)

	store := &mockStore{}
	conf := &devicerepository.Config{
		Store:         store,
		PhotosBaseURL: "https://assets/",
	}
	dr, err := devicerepository.New(c, conf)
	test.Must(dr, err)

	componenttest.StartComponent(t, c)
	defer c.Close()

	cc := dr.LoopbackConn()
	cl := ttnpb.NewDeviceRepositoryClient(cc)

	// conf is device repository config
	// store is device repository store
	// cl is device repository client

	// TODO: tests
}
