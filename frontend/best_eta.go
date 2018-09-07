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

package frontend

import (
	"context"
	"errors"
	"math"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/kelda-inc/hotrod-base/config"
	"github.com/kelda-inc/hotrod-base/pkg/pool"
	"github.com/kelda-inc/hotrod-customer/customer"
	"github.com/kelda-inc/hotrod-driver/driver"
	"github.com/kelda-inc/hotrod-route/route"
)

type bestETA struct {
	customer customer.Interface
	driver   driver.Interface
	route    route.Interface
	pool     *pool.Pool
}

// Response contains ETA for a trip.
type Response struct {
	Driver string
	DriverLocation string
	CustomerLocation string
	ETA    time.Duration
}

func newBestETA() *bestETA {
	return &bestETA{
		customer: customer.NewClient(),
		driver:   driver.NewClient(),
		route:    route.NewClient(),
		pool:     pool.New(config.RouteWorkerPoolSize),
	}
}

func (eta *bestETA) Get(ctx context.Context, customerID string) (*Response, error) {
	customer, err := eta.customer.Get(ctx, customerID)
	if err != nil {
		return nil, err
	}

	drivers, err := eta.driver.FindNearest(ctx, customer.Location)
	if err != nil {
		return nil, err
	}

	results := eta.getRoutes(ctx, customer, drivers)

	bestETAResp := &Response{ETA: math.MaxInt64}
	for _, result := range results {
		if result.err != nil {
			return nil, err
		}
		if result.route.ETA < bestETAResp.ETA {
			bestETAResp.ETA = result.route.ETA
			bestETAResp.Driver = result.driver
			bestETAResp.DriverLocation = result.route.Dropoff
			bestETAResp.CustomerLocation = result.route.Pickup
		}
	}
	if bestETAResp.Driver == "" {
		return nil, errors.New("No routes found")
	}

	return bestETAResp, nil
}

type routeResult struct {
	driver string
	route  *route.Route
	err    error
}

// getRoutes calls Route service for each (customer, driver) pair
func (eta *bestETA) getRoutes(ctx context.Context, customer *customer.Customer, drivers []driver.Driver) []routeResult {
	results := make([]routeResult, 0, len(drivers))
	wg := sync.WaitGroup{}
	routesLock := sync.Mutex{}
	for _, dd := range drivers {
		wg.Add(1)
		driver := dd // capture loop var
		// Use worker pool to (potentially) execute requests in parallel
		eta.pool.Execute(func() {
			route, err := eta.route.FindRoute(ctx, driver.Location, customer.Location)
			routesLock.Lock()
			results = append(results, routeResult{
				driver: driver.DriverID,
				route:  route,
				err:    err,
			})
			routesLock.Unlock()
			wg.Done()
		})
	}
	wg.Wait()
	return results
}

func (eta *bestETA) Transact(to, from string, amount int) error {
	time.Sleep(time.Second) // sleep a little to allow page to render

	ctx := context.TODO()

	eta.driver.Lock(ctx, to)
	eta.driver.Lock(ctx, from)

	err := eta.customer.Transfer(ctx, to, from, amount)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	eta.driver.Unlock(ctx, to)
	eta.driver.Unlock(ctx, from)

	return nil
}