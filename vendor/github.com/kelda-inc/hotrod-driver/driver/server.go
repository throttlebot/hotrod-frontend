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

package driver

import (
	log "github.com/sirupsen/logrus"
	"github.com/uber/tchannel-go"
	"github.com/uber/tchannel-go/thrift"

	"github.com/kelda-inc/hotrod-driver/driver/thrift-gen/driver"
	"time"
)

// Server implements jaeger-demo-frontend service
type Server struct {
	hostPort string
	ch       *tchannel.Channel
	server   *thrift.Server
	redis    *Redis
}

// NewServer creates a new driver.Server
func NewServer(hostPort string) *Server {
	channelOpts := &tchannel.ChannelOptions{
		//Tracer: tracer,
	}
	ch, err := tchannel.NewChannel("driver", channelOpts)
	if err != nil {
		log.WithError(err).Fatal("Cannot create TChannel")
	}
	server := thrift.NewServer(ch)

	return &Server{
		hostPort: hostPort,
		ch:       ch,
		server:   server,
		redis:    newRedis(),
	}
}

// Run starts the Driver server
func (s *Server) Run() error {

	s.server.Register(driver.NewTChanDriverServer(s))

	if err := s.ch.ListenAndServe(s.hostPort); err != nil {
		log.WithError(err).Fatal("Unable to start tchannel server")
	}

	peerInfo := s.ch.PeerInfo()
	log.WithField("hostPort", peerInfo.HostPort).Info("TChannel listening")

	// Run must block, but TChannel's ListenAndServe runs in the background, so block indefinitely
	select {}
}

// FindNearest implements Thrift interface TChanDriver
func (s *Server) FindNearest(ctx thrift.Context, location string) ([]*driver.DriverLocation, error) {
	log.WithField("location", location).Info("Searching for nearby drivers")
	driverIDs, err := s.redis.FindDriverIDs(ctx, location)
	if err != nil {
		log.WithError(err).Error("Failed to list drivers")
		return nil, err
	}

	retMe := make([]*driver.DriverLocation, len(driverIDs))
	for i, driverID := range driverIDs {
		var drv Driver
		var err error
		for i := 0; i < 3; i++ {
			drv, err = s.redis.GetDriver(ctx, driverID)
			if err == nil {
				break
			}
			log.WithError(err).WithField("retry_no", i+1).Error("Retrying GetDriver after error")
		}
		if err != nil {
			log.WithError(err).Error("Failed to get driver after 3 attempts")
			return nil, err
		}
		retMe[i] = &driver.DriverLocation{
			DriverID: drv.DriverID,
			Location: drv.Location,
		}
	}
	log.WithField("num_drivers", len(retMe)).Info("Search successful")
	return retMe, nil
}

// Lock uses redis to implement a lock on an account
func (s *Server) Lock(ctx thrift.Context, id string) (*driver.Result_, error) {
	log.WithField(id, id).Info("Attempting to secure lock")
	for !s.redis.AttemptLock(ctx, id) {
		time.Sleep(time.Millisecond * 100)
	}
	log.WithField(id, id).Info("Secured lock")
	return &driver.Result_{}, nil
}

// Lock uses redis to implement a lock on an account
func (s *Server) Unlock(ctx thrift.Context, id string) (*driver.Result_, error) {
	log.WithField(id, id).Info("Releasing lock")
	s.redis.Unlock(id)
	log.WithField(id, id).Info("Released lock")
	return &driver.Result_{}, nil
}