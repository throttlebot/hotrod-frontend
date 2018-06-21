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

package customer

import (
	"context"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"gitlab.com/will.wang1/hotrod-base/pkg/tracing"
)

// Client is a remote client that implements customer.Interface
type Client struct {
	client *tracing.HTTPClient
}

// NewClient creates a new customer.Client
func NewClient() *Client {
	return &Client{
		client: &tracing.HTTPClient{
			Client: http.DefaultClient,
		},
	}
}

// Get implements customer.Interface#Get as an RPC
func (c *Client) Get(ctx context.Context, customerID string) (*Customer, error) {
	log.WithField("customer_id", customerID).Info("Getting customer")

	clientIP := "hotrod-customer:8081"

	url := fmt.Sprintf("http://" + clientIP + "/customer?customer=%s", customerID)
	var customer Customer
	if err := c.client.GetJSON(ctx, "/customer", url, &customer); err != nil {
		return nil, err
	}
	return &customer, nil
}
