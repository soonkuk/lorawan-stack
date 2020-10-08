// Copyright © 2020 The Things Industries B.V.

package storage

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages/storage/provider"
	"go.thethings.network/lorawan-stack/v3/pkg/auth/rights"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"google.golang.org/grpc"
)

// RegisterServices implements packages.ApplicationPackageHandler.
func (p *storage) RegisterServices(s *grpc.Server) {
	ttnpb.RegisterApplicationUpStorageServer(s, p)
}

// RegisterHandlers implements packages.ApplicationPackageHandler.
func (p *storage) RegisterHandlers(s *runtime.ServeMux, conn *grpc.ClientConn) {
	ttnpb.RegisterApplicationUpStorageHandler(p.ctx, s, conn)
}

// GetStoredApplicationUp implements ttnpb.ApplicationUpStorage.
func (p *storage) GetStoredApplicationUp(req *ttnpb.GetStoredApplicationUpRequest, stream ttnpb.ApplicationUpStorage_GetStoredApplicationUpServer) error {
	ids, err := applicationIDs(req)
	if err != nil {
		return err
	}
	ctx := stream.Context()
	if err := rights.RequireApplication(ctx, *ids, ttnpb.RIGHT_APPLICATION_TRAFFIC_READ); err != nil {
		return err
	}
	return p.provider.Range(ctx, provider.QueryFromRequest(req), func(up *ttnpb.ApplicationUp) error {
		return stream.Send(up)
	})
}

func applicationIDs(req *ttnpb.GetStoredApplicationUpRequest) (*ttnpb.ApplicationIdentifiers, error) {
	if req.EndDeviceIDs != nil {
		if req.EndDeviceIDs.ApplicationIdentifiers.IsZero() && req.ApplicationIDs != nil {
			req.EndDeviceIDs.ApplicationIdentifiers = *req.ApplicationIDs
		}
		return &req.EndDeviceIDs.ApplicationIdentifiers, nil
	}
	if req.ApplicationIDs == nil {
		return nil, errNoAppID.New()
	}
	return req.ApplicationIDs, nil
}
