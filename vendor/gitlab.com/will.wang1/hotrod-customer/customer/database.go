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
	"database/sql"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"os"
	"fmt"
)

// database simulates Customer repository implemented on top of an SQL database
type database struct {
	*sql.DB
}

func newDatabase() (*database, error) {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASS")
	url := os.Getenv("POSTGRES_URL")
	connectStr := fmt.Sprintf("postgres://%s:%s@%s/customers?sslmode=disable", user, password, url)
	db, err := sql.Open("postgres", connectStr)
	return &database{db}, err
}

func (d *database) Get(ctx context.Context, customerID string) (*Customer, error) {
	log.WithField("customer_id", customerID).Info("Loading customer")

	customer := Customer{ID: customerID}
	err := d.QueryRow("SELECT name, location FROM customers WHERE id = $1", customerID).Scan(
		&customer.Name, &customer.Location)
	return &customer, err
}
