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
	"context"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"

	"github.com/kelda-inc/hotrod-base/pkg/delay"
	"github.com/kelda-inc/hotrod-base/config"

	"os"
	"time"
)

// Redis is a simulator of remote Redis cache
type Redis struct {
	*redis.Client
}

func newRedis() *Redis {
	redisURL := os.Getenv("REDIS_URL")
	redisPass := os.Getenv("REDIS_PASS")

	client := redis.NewClient(&redis.Options{
		Addr: redisURL,
		Password: redisPass,
	})
	return &Redis{client}

}

// FindDriverIDs finds IDs of drivers who are near the location.
func (r *Redis) FindDriverIDs(ctx context.Context, location string) ([]string, error) {
	// simulate RPC delay
	delay.Sleep(config.RedisFindDelay, config.RedisFindDelayStdDev)
	return r.Keys("T7*").Result()
}

// GetDriver returns driver and the current car location
func (r *Redis) GetDriver(ctx context.Context, driverID string) (Driver, error) {
	// simulate RPC delay
	delay.Sleep(config.RedisGetDelay, config.RedisGetDelayStdDev)
	driver, err := r.Get(driverID).Result()
	if err != nil {
		log.WithField("driver_id", driverID).WithError(err).Error("failed to get driver")
		return Driver{}, err
	}

	return Driver{
		DriverID: driverID,
		Location: driver,
	}, nil
}

// AttemptLock calls SETNX and returns if the lock was acquired
func (r *Redis) AttemptLock(ctx context.Context, id string) bool {
	// simulate RPC delay
	time.Sleep(time.Millisecond * 1500)
	return r.SetNX("lock-" + id, 1, time.Minute).Val()
}

func (r *Redis) Unlock(id string) {
	// simulate RPC delay
	time.Sleep(time.Millisecond * 1500)
	r.Del("lock-" + id)
}