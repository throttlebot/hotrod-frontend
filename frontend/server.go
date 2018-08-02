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
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/kelda-inc/hotrod-base/pkg/httperr"
	"github.com/kelda-inc/hotrod-base/pkg/tracing"
	"github.com/kelda-inc/hotrod-customer/customer"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"
)

// Server implements jaeger-demo-frontend service
// New comment
type Server struct {
	customerClient customer.Interface
	hostPort       string
	bestETA        *bestETA
}

// NewServer creates a new frontend.Server
func NewServer(hostPort string) *Server {
	return &Server{
		customerClient: customer.NewClient(),
		hostPort:       hostPort,
		bestETA:        newBestETA(),
	}
}

var httpReqs = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "frontend_http_requests_total",
		Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
	},
	[]string{"code", "method", "path"},
)

// Run starts the frontend server
func (s *Server) Run() error {
	mux := s.createServeMux()

	prometheus.MustRegister(httpReqs)

	log.WithField("address", "http://"+s.hostPort).Info("Starting")
	return http.ListenAndServe(s.hostPort, mux)
}

func (s *Server) createServeMux() http.Handler {
	mux := tracing.NewServeMux()
	mux.Handle("/customers", http.HandlerFunc(s.customers))
	mux.Handle("/dispatch", http.HandlerFunc(s.dispatch))
	mux.Handle("/refund", http.HandlerFunc(s.refund))
	mux.Handle("/metrics", promhttp.Handler())
	return mux
}

func (s *Server) customers(w http.ResponseWriter, r *http.Request) {
	log.WithField("method", r.Method).WithField("url", r.URL).Info("HTTP")

	customers, err := s.customerClient.ListCustomerPublicInfo(r.Context())
	if httperr.HandleError(w, err, http.StatusInternalServerError) {
		log.WithError(err).Error("Failed to query customers")
		httpReqs.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.Method, r.URL.Path).Inc()
		return
	}

	data, err := json.Marshal(customers)

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

	httpReqs.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method, r.URL.Path).Inc()


}

func (s *Server) dispatch(w http.ResponseWriter, r *http.Request) {
	log.WithField("method", r.Method).WithField("url", r.URL).Info("HTTP request received")
	if err := r.ParseForm(); httperr.HandleError(w, err, http.StatusBadRequest) {
		log.WithError(err).Error("bad request")
		httpReqs.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.Method, r.URL.Path).Inc()
		return
	}

	customerID := r.Form.Get("customer")
	if customerID == "" {
		http.Error(w, "Missing required 'customer' parameter", http.StatusBadRequest)
		httpReqs.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.Method, r.URL.Path).Inc()
		return
	}

	ctx := r.Context()

	// TODO distinguish between user errors (such as invalid customer ID) and server failures
	response, err := s.bestETA.Get(ctx, customerID)
	if httperr.HandleError(w, err, http.StatusInternalServerError) {
		log.WithError(err).Error("request failed")
		httpReqs.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.Method, r.URL.Path).Inc()
		return
	}

	data, err := json.Marshal(response)
	if httperr.HandleError(w, err, http.StatusInternalServerError) {
		log.WithError(err).Error("cannot marshal response")
		httpReqs.WithLabelValues(strconv.Itoa(http.StatusInternalServerError), r.Method, r.URL.Path)
		return
	}


	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

	httpReqs.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method, r.URL.Path).Inc()

	go s.bestETA.Transact(response.Driver, customerID, 5)


}

func (s *Server) refund(w http.ResponseWriter, r *http.Request) {
	log.WithField("method", r.Method).WithField("url", r.URL).Info("HTTP request received")
	if err := r.ParseForm(); httperr.HandleError(w, err, http.StatusBadRequest) {
		log.WithError(err).Error("bad request")
		httpReqs.WithLabelValues(strconv.Itoa(http.StatusBadRequest), r.Method, r.URL.Path).Inc()
		return
	}
	customerID := r.Form.Get("customer")
	if customerID == "" {
		http.Error(w, "Missing required 'customer' parameter", http.StatusBadRequest)
		return
	}
	driverID := r.Form.Get("driver")
	if driverID == "" {
		http.Error(w, "Missing required 'driver' parameter", http.StatusBadRequest)
		return
	}

	if err := s.bestETA.Transact(customerID, driverID, 5); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("success"))

}

