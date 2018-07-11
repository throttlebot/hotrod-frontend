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
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/kelda-inc/hotrod-base/pkg/tracing"
)

// Client is a remote client that implements customer.Interface
type Client struct {
	client 	*tracing.HTTPClient
	address string
}

// NewClient creates a new customer.Client
func NewClient() *Client {
	customerHost := os.Getenv("HOTROD_CUSTOMER_HOST")
	if customerHost == "" {
		customerHost = "hotrod-customer"
	}

	return &Client{
		client: &tracing.HTTPClient{
			Client: http.DefaultClient,
		},
		address: "http://" + customerHost + ":8081",
	}
}

// Get implements customer.Interface#Get as an RPC
func (c *Client) Get(ctx context.Context, customerID string) (*Customer, error) {
	log.WithField("customer_id", customerID).Info("Getting customer")

	url := fmt.Sprintf(c.address + "/customer?customer=%s", customerID)
	var customer Customer
	if err := c.client.GetJSON(ctx, "/customer", url, &customer); err != nil {
		return nil, err
	}
	return &customer, nil
}

func (c *Client) ListCustomerPublicInfo(ctx context.Context) ([]Customer, error) {
	log.Info("Getting all customers")
	url := c.address + "/list"
	var customers []Customer
	if err := c.client.GetJSON(ctx, "/customer", url, &customers); err != nil {
		return nil, err
	}
	return customers, nil
}
