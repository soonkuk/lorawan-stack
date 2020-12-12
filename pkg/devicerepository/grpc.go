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
	"strconv"
	"strings"

	clusterauth "go.thethings.network/lorawan-stack/v3/pkg/auth/cluster"
	"go.thethings.network/lorawan-stack/v3/pkg/devicerepository/store"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// withDefaultModelFields appends default ttnpb.EndDeviceModel fields.
func withDefaultModelFields(paths []string) []string {
	return ttnpb.AddFields(paths, "brand_id", "model_id")
}

// withDefaultBrandFields appends default ttnpb.EndDeviceBrand paths.
func withDefaultBrandFields(paths []string) []string {
	return ttnpb.AddFields(paths, "brand_id")
}

func (dr *DeviceRepository) assetURL(brandID, path string) string {
	if path == "" || dr.config.AssetsBaseURL == "" {
		return path
	}
	return strings.TrimRight(dr.config.AssetsBaseURL, "/") + "/vendor/" + brandID + "/" + path
}

// ListBrands implements the ttnpb.DeviceRepositoryServer interface.
func (dr *DeviceRepository) ListBrands(ctx context.Context, request *ttnpb.ListEndDeviceBrandsRequest) (*ttnpb.ListEndDeviceBrandsResponse, error) {
	if dr.config.RequireAuth {
		// TODO: require any application rights here.
	}
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
	for _, brand := range response.Brands {
		brand.Logo = dr.assetURL(brand.BrandID, brand.Logo)
	}
	grpc.SetHeader(ctx, metadata.Pairs("x-total-count", strconv.FormatUint(uint64(response.Total), 10)))
	return &ttnpb.ListEndDeviceBrandsResponse{
		Brands: response.Brands,
	}, nil
}

var (
	errBrandNotFound = errors.DefineNotFound("brand_not_found", "brand `{brand_id}` not found")
)

// GetBrand implements the ttnpb.DeviceRepositoryServer interface.
func (dr *DeviceRepository) GetBrand(ctx context.Context, request *ttnpb.GetEndDeviceBrandRequest) (*ttnpb.EndDeviceBrand, error) {
	if dr.config.RequireAuth {
		// TODO: require any application rights here.
	}
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
	brand := response.Brands[0]
	brand.Logo = dr.assetURL(brand.BrandID, brand.Logo)
	return brand, nil
}

// ListModels implements the ttnpb.DeviceRepositoryServer interface.
func (dr *DeviceRepository) ListModels(ctx context.Context, request *ttnpb.ListEndDeviceModelsRequest) (*ttnpb.ListEndDeviceModelsResponse, error) {
	if dr.config.RequireAuth {
		// TODO: require any application rights here.
	}
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

	for _, model := range response.Models {
		if photos := model.Photos; photos != nil {
			photos.Main = dr.assetURL(model.BrandID, photos.Main)
			for idx, photo := range photos.Other {
				photos.Other[idx] = dr.assetURL(model.BrandID, photo)
			}
		}
	}

	grpc.SetHeader(ctx, metadata.Pairs("x-total-count", strconv.FormatUint(uint64(response.Total), 10)))
	return &ttnpb.ListEndDeviceModelsResponse{
		Models: response.Models,
	}, nil
}

var (
	errModelNotFound = errors.DefineNotFound("model_not_found", "model `{brand_id}/{model_id}` not found")
)

// GetModel implements the ttnpb.DeviceRepositoryServer interface.
func (dr *DeviceRepository) GetModel(ctx context.Context, request *ttnpb.GetEndDeviceModelRequest) (*ttnpb.EndDeviceModel, error) {
	if dr.config.RequireAuth {
		// TODO: require any application rights here.
	}
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
	model := response.Models[0]
	if photos := model.Photos; photos != nil {
		photos.Main = dr.assetURL(model.BrandID, photos.Main)
		for idx, photo := range photos.Other {
			photos.Other[idx] = dr.assetURL(model.BrandID, photo)
		}
	}

	return model, nil
}

// GetTemplate implements the ttnpb.DeviceRepositoryServer interface.
func (dr *DeviceRepository) GetTemplate(ctx context.Context, ids *ttnpb.EndDeviceVersionIdentifiers) (*ttnpb.EndDeviceTemplate, error) {
	if dr.config.RequireAuth {
		// TODO: require any application rights here.
	}
	return dr.store.GetTemplate(ids)
}

// GetUplinkDecoder implements the ttnpb.DeviceRepositoryServer interface.
func (dr *DeviceRepository) GetUplinkDecoder(ctx context.Context, ids *ttnpb.EndDeviceVersionIdentifiers) (*ttnpb.MessagePayloadFormatter, error) {
	if dr.config.RequireAuth {
		// TODO: require any application rights here.
		if err := clusterauth.Authorized(ctx); err != nil {
			return nil, err
		}
	}
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
	if dr.config.RequireAuth {
		// TODO: require any application rights here.
		if err := clusterauth.Authorized(ctx); err != nil {
			return nil, err
		}
	}
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
	if dr.config.RequireAuth {
		// TODO: require any application rights here.
		if err := clusterauth.Authorized(ctx); err != nil {
			return nil, err
		}
	}
	s, err := dr.store.GetDownlinkEncoder(ids)
	if err != nil {
		return nil, err
	}
	return &ttnpb.MessagePayloadFormatter{
		Formatter:          ttnpb.PayloadFormatter_FORMATTER_JAVASCRIPT,
		FormatterParameter: s,
	}, nil
}
