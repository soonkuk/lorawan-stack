// Copyright © 2019 The Things Network Foundation, The Things Industries B.V.
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

package mock

import (
	"context"
	"crypto/tls"
	"crypto/x509"

	routingpb "go.packetbroker.org/api/routing"
	packetbroker "go.packetbroker.org/api/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// PBDataPlane is a mock Packet Broker Data Plane.
type PBDataPlane struct {
	*grpc.Server
	ForwarderUp     chan *packetbroker.RoutedUplinkMessage
	ForwarderDown   chan *packetbroker.RoutedDownlinkMessage
	HomeNetworkDown chan *packetbroker.RoutedDownlinkMessage
	HomeNetworkUp   chan *packetbroker.RoutedUplinkMessage
}

// NewPBDataPlane instantiates a new mock Packet Broker Data Plane.
func NewPBDataPlane(cert tls.Certificate, clientCAs *x509.CertPool) *PBDataPlane {
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCAs,
	})
	dp := &PBDataPlane{
		Server: grpc.NewServer(
			grpc.Creds(creds),
		),
		ForwarderUp:     make(chan *packetbroker.RoutedUplinkMessage),
		ForwarderDown:   make(chan *packetbroker.RoutedDownlinkMessage),
		HomeNetworkDown: make(chan *packetbroker.RoutedDownlinkMessage),
		HomeNetworkUp:   make(chan *packetbroker.RoutedUplinkMessage),
	}
	routingpb.RegisterForwarderDataServer(dp.Server, &routerForwarderServer{
		upCh:   dp.ForwarderUp,
		downCh: dp.ForwarderDown,
	})
	routingpb.RegisterHomeNetworkDataServer(dp.Server, &routerHomeNetworkServer{
		downCh: dp.HomeNetworkDown,
		upCh:   dp.HomeNetworkUp,
	})
	return dp
}

type routerForwarderServer struct {
	routingpb.UnimplementedForwarderDataServer
	upCh   chan *packetbroker.RoutedUplinkMessage
	downCh chan *packetbroker.RoutedDownlinkMessage
}

func (s *routerForwarderServer) Publish(ctx context.Context, req *routingpb.PublishUplinkMessageRequest) (*routingpb.PublishUplinkMessageResponse, error) {
	s.upCh <- &packetbroker.RoutedUplinkMessage{
		Message: req.Message,
	}
	return &routingpb.PublishUplinkMessageResponse{
		Id: "test",
	}, nil
}

func (s *routerForwarderServer) Subscribe(req *routingpb.SubscribeForwarderRequest, res routingpb.ForwarderData_SubscribeServer) error {
	for {
		select {
		case <-res.Context().Done():
			return nil
		case msg := <-s.downCh:
			if err := res.Send(msg); err != nil {
				return err
			}
		}
	}
}

type routerHomeNetworkServer struct {
	routingpb.UnimplementedHomeNetworkDataServer
	downCh chan *packetbroker.RoutedDownlinkMessage
	upCh   chan *packetbroker.RoutedUplinkMessage
}

func (s *routerHomeNetworkServer) Publish(ctx context.Context, req *routingpb.PublishDownlinkMessageRequest) (*routingpb.PublishDownlinkMessageResponse, error) {
	down := &packetbroker.RoutedDownlinkMessage{
		ForwarderNetId:     req.ForwarderNetId,
		ForwarderTenantId:  req.ForwarderTenantId,
		ForwarderClusterId: req.ForwarderClusterId,
		Message:            req.Message,
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case s.downCh <- down:
	}
	return &routingpb.PublishDownlinkMessageResponse{
		Id: "test",
	}, nil
}

func (s *routerHomeNetworkServer) Subscribe(req *routingpb.SubscribeHomeNetworkRequest, res routingpb.HomeNetworkData_SubscribeServer) error {
	for {
		select {
		case <-res.Context().Done():
			return nil
		case msg := <-s.upCh:
			if err := res.Send(msg); err != nil {
				return err
			}
		}
	}
}
