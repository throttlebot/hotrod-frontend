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
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/kelda-inc/hotrod-base/pkg/httperr"
	"github.com/kelda-inc/hotrod-base/pkg/tracing"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"fmt"
)

// Server implements Customer service
type Server struct {
	hostPort string
	database *database
}

// NewServer creates a new customer.Server
func NewServer(hostPort string) (*Server, error) {
	db, err := newDatabase()
	if err != nil {
		return nil, err
	}
	return &Server{
		hostPort: hostPort,
		database: db,
	}, nil
}

// Run starts the Customer server
func (s *Server) Run() error {
	mux := s.createServeMux()
	log.WithField("address", "http://"+s.hostPort).Info("Starting")
	return http.ListenAndServe(s.hostPort, mux)
}

func (s *Server) createServeMux() http.Handler {
	mux := tracing.NewServeMux()
	mux.Handle("/customer", http.HandlerFunc(s.customer))
	mux.Handle("/transfer", http.HandlerFunc(s.transfer))
	mux.Handle("/list", http.HandlerFunc(s.listCustomers))
	mux.Handle("/metrics", promhttp.Handler())
	return mux
}

func (s *Server) customer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.WithField("method", r.Method).WithField("url", r.URL).Info("HTTP request received")
	if err := r.ParseForm(); httperr.HandleError(w, err, http.StatusBadRequest) {
		log.WithError(err).Error("bad request")
		return
	}

	customerID := r.Form.Get("customer")
	if customerID == "" {
		http.Error(w, "Missing required 'customer' parameter", http.StatusBadRequest)
		return
	}

	response, err := s.database.Get(ctx, customerID)
	if httperr.HandleError(w, err, http.StatusInternalServerError) {
		log.WithError(err).Error("request failed")
		return
	}

	data, err := json.Marshal(response)
	if httperr.HandleError(w, err, http.StatusInternalServerError) {
		log.WithError(err).Error("cannot marshal response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (s *Server) listCustomers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.WithField("method", r.Method).WithField("url", r.URL).Info("HTTP request received")

	response, err := s.database.List(ctx)
	if httperr.HandleError(w, err, http.StatusInternalServerError) {
		log.WithError(err).Error("request failed")
		return
	}

	data, err := json.Marshal(response)
	if httperr.HandleError(w, err, http.StatusInternalServerError) {
		log.WithError(err).Error("cannot marshal response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (s *Server) transfer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	toID := r.Form.Get("to")
	fromID := r.Form.Get("from")
	amount := r.Form.Get("amount")
	if toID == "" || fromID == "" || amount == "" {
		http.Error(w, fmt.Sprintf("Missing field: to=%s, from=%s, amount=%s", toID, fromID, amount), http.StatusBadRequest)
		return
	}

	err := s.database.Transfer(ctx, toID, fromID, amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("success"))
}