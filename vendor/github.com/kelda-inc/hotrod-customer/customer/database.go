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
	"strconv"
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



func (d *database) List(ctx context.Context ) ([]Customer, error) {
	log.Info("Loading all customers")
	rows, err := d.Query("SELECT name, id FROM customers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers := make([]Customer, 0)
	for rows.Next() {
		var id, name string
		rows.Scan(&name, &id)
		customers = append(customers, Customer{
			Name: name,
			ID: id,
		})
	}

	return customers, rows.Err()
}


func(d *database) Transfer(ctx context.Context, to, from, amount string) error {
	i, err := strconv.Atoi(amount)
	if err != nil {
		return err
	}

	// Check if from account exists and has enough money
	account := Account{ID: from}
	err = d.QueryRow("SELECT balance FROM accounts WHERE id = $1", from).Scan(
		&account.Balance)
	if err == sql.ErrNoRows {
		return fmt.Errorf("account %s does not exist")
	} else if err != nil {
		return err
	} else if account.Balance - i <= 0 {
		return fmt.Errorf("account %s does not have anough money: %d",account.ID, account.Balance)
	}

	// Checks if to account exists and if not init the new account
	account = Account{ID: to}
	err = d.QueryRow("SELECT balance FROM accounts WHERE id = $1", to).Scan(
		&account.Balance)
	if err == sql.ErrNoRows {
		_, err = d.Exec("INSERT INTO accounts VALUES ($1, $2)", to, 0)
	}
	if err != nil {
		return err
	}

	// Transfer money (assume that this does NOT fail!)
	_, err = d.Exec("UPDATE accounts SET balance = balance - $1 WHERE id = $2", amount, from)
	if err != nil {
		return err
	}
	_, err = d.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, to)
	return err
}