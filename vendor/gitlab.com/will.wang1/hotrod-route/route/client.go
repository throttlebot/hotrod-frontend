// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package route

import (
	"context"
	"net/http"
	"net/url"
	"os"

	log "github.com/sirupsen/logrus"

	"gitlab.com/will.wang1/hotrod-base/pkg/tracing"
)

// Client is a remote client that implements route.Interface
type Client struct {
	client    *tracing.HTTPClient
	routeHost string
}

// NewClient creates a new route.Client
func NewClient() *Client {
	routeHost := os.Getenv("HOTROD_ROUTE_HOST")
	if routeHost == "" {
		routeHost = "hotrod-route"
	}
	routeHost += ":8083"
	return &Client{
		client: &tracing.HTTPClient{
			Client: http.DefaultClient,
		},
		routeHost: routeHost,
	}
}

// FindRoute implements route.Interface#FindRoute as an RPC
func (c *Client) FindRoute(ctx context.Context, pickup, dropoff string) (*Route, error) {
	log.WithField("pickup", pickup).WithField("dropoff", dropoff).Info("Finding route")

	v := url.Values{}
	v.Set("pickup", pickup)
	v.Set("dropoff", dropoff)

	url := "http://" + c.routeHost + "/route?" + v.Encode()

	var route Route
	if err := c.client.GetJSON(ctx, "/route", url, &route); err != nil {
		log.WithError(err).Error("Error getting route")
		return nil, err
	}
	return &route, nil
}
