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
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

// withDefaultModelFields appends default ttnpb.EndDeviceModel fields.
func withDefaultModelFields(paths []string) []string {
	return ttnpb.AddFields(paths, "brand_id", "model_id")
}

// withDefaultBrandFields appends default ttnpb.EndDeviceBrand paths.
func withDefaultBrandFields(paths []string) []string {
	return ttnpb.AddFields(paths, "brand_id")
}

// ListBrands implements the ttnpb.DeviceRepositoryServer interface.
func (dr *DeviceRepository) ListBrands(ctx context.Context, request *ttnpb.ListEndDeviceBrandsRequest) (*ttnpb.ListEndDeviceBrandsResponse, error) {
	if request.Limit == 0 {
		request.Limit = 1000
	}
	response, err := dr.store.GetBrands(store.GetBrandsRequest{
		Limit:   request.Limit,
		Page:    request.Page,
		OrderBy: request.OrderBy,
		Paths:   withDefaultBrandFields(request.FieldMask.Paths),
		Search:  request.Search,
	})
	if err != nil {
		return nil, err
	}
	return &ttnpb.ListEndDeviceBrandsResponse{
		Brands: response.Brands,
		Count:  response.Count,
		Offset: response.Offset,
		Total:  response.Total,
	}, nil
}

var (
	errBrandNotFound = errors.DefineNotFound("brand_not_found", "brand `{brand_id}` not found")
)

// GetBrand implements the ttnpb.DeviceRepositoryServer interface.
func (dr *DeviceRepository) GetBrand(ctx context.Context, request *ttnpb.GetEndDeviceBrandRequest) (*ttnpb.EndDeviceBrand, error) {
	response, err := dr.store.GetBrands(store.GetBrandsRequest{
		BrandID: request.BrandID,
		Paths:   withDefaultBrandFields(request.FieldMask.Paths),
		Limit:   1,
	})
	if err != nil {
		return nil, err
	}
	if len(response.Brands) == 0 {
		return nil, errBrandNotFound.WithAttributes("brand_id", request.BrandID)
	}
	return response.Brands[0], nil
}

// ListModels implements the ttnpb.DeviceRepositoryServer interface.
func (dr *DeviceRepository) ListModels(ctx context.Context, request *ttnpb.ListEndDeviceModelsRequest) (*ttnpb.ListEndDeviceModelsResponse, error) {
	if request.Limit == 0 {
		request.Limit = 1000
	}
	response, err := dr.store.GetModels(store.GetModelsRequest{
		BrandID: request.BrandID,
		Limit:   request.Limit,
		Page:    request.Page,
		Paths:   withDefaultModelFields(request.FieldMask.Paths),
		Search:  request.Search,
		OrderBy: request.OrderBy,
	})
	if err != nil {
		return nil, err
	}
	return &ttnpb.ListEndDeviceModelsResponse{
		Models: response.Models,
		Count:  response.Count,
		Offset: response.Offset,
		Total:  response.Total,
	}, nil
}

var (
	errModelNotFound = errors.DefineNotFound("model_not_found", "model `{brand_id}/{model_id}` not found")
)

// GetModel implements the ttnpb.DeviceRepositoryServer interface.
func (dr *DeviceRepository) GetModel(ctx context.Context, request *ttnpb.GetEndDeviceModelRequest) (*ttnpb.EndDeviceModel, error) {
	response, err := dr.store.GetModels(store.GetModelsRequest{
		BrandID: request.BrandID,
		ModelID: request.ModelID,
		Limit:   1,
		Paths:   withDefaultModelFields(request.FieldMask.Paths),
	})
	if err != nil {
		return nil, err
	}
	if len(response.Models) == 0 {
		return nil, errModelNotFound.WithAttributes("brand_id", request.BrandID, "model_id", request.ModelID)
	}
	return response.Models[0], nil
}

// GetTemplate implements the ttnpb.DeviceRepositoryServer interface.
func (dr *DeviceRepository) GetTemplate(ctx context.Context, ids *ttnpb.EndDeviceVersionIdentifiers) (*ttnpb.EndDeviceTemplate, error) {
	return dr.store.GetTemplate(ids)
}

// GetUplinkDecoder implements the ttnpb.DeviceRepositoryServer interface.
func (dr *DeviceRepository) GetUplinkDecoder(ctx context.Context, ids *ttnpb.EndDeviceVersionIdentifiers) (*ttnpb.MessagePayloadFormatter, error) {
	s, err := dr.store.GetUplinkDecoder(ids)
	if err != nil {
		return nil, err
	}
	return &ttnpb.MessagePayloadFormatter{
		Formatter:          ttnpb.PayloadFormatter_FORMATTER_JAVASCRIPT,
		FormatterParameter: s,
	}, nil
}

// GetDownlinkDecoder implements the ttnpb.DeviceRepositoryServer interface.
func (dr *DeviceRepository) GetDownlinkDecoder(ctx context.Context, ids *ttnpb.EndDeviceVersionIdentifiers) (*ttnpb.MessagePayloadFormatter, error) {
	s, err := dr.store.GetDownlinkDecoder(ids)
	if err != nil {
		return nil, err
	}
	return &ttnpb.MessagePayloadFormatter{
		Formatter:          ttnpb.PayloadFormatter_FORMATTER_JAVASCRIPT,
		FormatterParameter: s,
	}, nil
}

// GetDownlinkEncoder implements the ttnpb.DeviceRepositoryServer interface.
func (dr *DeviceRepository) GetDownlinkEncoder(ctx context.Context, ids *ttnpb.EndDeviceVersionIdentifiers) (*ttnpb.MessagePayloadFormatter, error) {
	s, err := dr.store.GetDownlinkEncoder(ids)
	if err != nil {
		return nil, err
	}
	return &ttnpb.MessagePayloadFormatter{
		Formatter:          ttnpb.PayloadFormatter_FORMATTER_JAVASCRIPT,
		FormatterParameter: s,
	}, nil
}
