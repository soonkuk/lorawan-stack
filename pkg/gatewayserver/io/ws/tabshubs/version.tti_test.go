// Copyright © 2020 The Things Industries B.V.

package tabshubs

import (
	"encoding/json"
	"testing"

	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
)

func TestType(t *testing.T) {
	a := assertions.New(t)
	msg := Version{
		Station:  "test",
		Firmware: "2.0.0",
		Package:  "test",
		Model:    "test",
		Protocol: 2,
	}

	data, err := json.Marshal(msg)
	a.So(err, should.BeNil)

	mt, err := Type(data)
	a.So(err, should.BeNil)
	a.So(mt, should.Equal, TypeUpstreamVersion)
}

func TestIsProduction(t *testing.T) {
	for _, tc := range []struct {
		Name             string
		Message          Version
		ExpectedResponse bool
	}{
		{
			Name:             "EmptyMessage",
			Message:          Version{},
			ExpectedResponse: false,
		},
		{
			Name: "EmptyMessage1",
			Message: Version{
				Features: "",
			},
			ExpectedResponse: false,
		},
		{
			Name: "NonProduction",
			Message: Version{
				Features: "gps rmtsh",
			},
			ExpectedResponse: false,
		},
		{
			Name: "Production",
			Message: Version{
				Features: "prod",
			},
			ExpectedResponse: true,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			a := assertions.New(t)
			a.So(tc.Message.IsProduction(), should.Equal, tc.ExpectedResponse)
		})
	}
}
