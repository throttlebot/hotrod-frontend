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
	"encoding/json"
	"expvar"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"gitlab.com/will.wang1/hotrod-base/config"
	"gitlab.com/will.wang1/hotrod-base/pkg/delay"
	"gitlab.com/will.wang1/hotrod-base/pkg/httperr"
	"gitlab.com/will.wang1/hotrod-base/pkg/tracing"
)

// Server implements Route service
type Server struct {
	hostPort string
}

// NewServer creates a new route.Server
func NewServer(hostPort string) *Server {
	return &Server{
		hostPort: hostPort,
	}
}

// Run starts the Route server
func (s *Server) Run() error {
	mux := s.createServeMux()
	log.WithField("address", "http://"+s.hostPort).Info("Starting")
	return http.ListenAndServe(s.hostPort, mux)
}

func (s *Server) createServeMux() http.Handler {
	mux := tracing.NewServeMux()
	mux.Handle("/route", http.HandlerFunc(s.route))
	mux.Handle("/debug/vars", expvar.Handler()) // expvar
	mux.Handle("/metrics", promhttp.Handler())  // Prometheus
	return mux
}

func (s *Server) route(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.WithField("method", r.Method).WithField("url", r.URL).Info("HTTP request received")
	if err := r.ParseForm(); httperr.HandleError(w, err, http.StatusBadRequest) {
		log.WithError(err).Error("bad request")
		return
	}

	pickup := r.Form.Get("pickup")
	if pickup == "" {
		http.Error(w, "Missing required 'pickup' parameter", http.StatusBadRequest)
		return
	}

	dropoff := r.Form.Get("dropoff")
	if dropoff == "" {
		http.Error(w, "Missing required 'dropoff' parameter", http.StatusBadRequest)
		return
	}

	response := computeRoute(ctx, pickup, dropoff)

	data, err := json.Marshal(response)
	if httperr.HandleError(w, err, http.StatusInternalServerError) {
		log.WithError(err).Error("cannot marshal response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func computeRoute(ctx context.Context, pickup, dropoff string) *Route {
	start := time.Now()
	defer func() {
		updateCalcStats(ctx, time.Since(start))
	}()

	// Simulate expensive calculation
	delay.Sleep(config.RouteCalcDelay, config.RouteCalcDelayStdDev)

	return &Route{
		Pickup:  pickup,
		Dropoff: dropoff,
		ETA:     computeEta(),
	}
}

func computeEta() time.Duration {
	seconds := int16(rand.Intn(32967))
	return time.Duration(seconds/60) * time.Minute // Returns a time.Duration object in minutes
}
