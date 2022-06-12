/*
 * Copyright 2022 Sue B.V.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package health

import (
	"context"
	"encoding/json"
	"github.com/suecodelabs/cnfuzz/src/log"
	"net/http"
	"time"
)

const (
	StatusKey       = "status"
	UnHealthyStatus = "fail"
	HealthyStatus   = "pass"
)

type ICheck interface {
	CheckHealth(context context.Context) Health
}

type Check struct {
	name  string
	check ICheck
}

type Checker struct {
	checkers []Check
}

func NewChecker() Checker {
	return Checker{
		checkers: make([]Check, 0),
	}
}

func (c *Checker) RegisterCheck(name string, check ICheck) {
	if check == nil {
		log.L().Warnf("failed to register %s health check, because it doesn't contain a check function", name)
		return
	}
	newCheck := Check{name: name, check: check}
	c.checkers = append(c.checkers, newCheck)
}

func (c *Checker) IsHealthy() Health {
	h := NewHealth(true)

	ctx, _ := context.WithTimeout(context.TODO(), time.Second*3)
	for _, checker := range c.checkers {
		health := checker.check.CheckHealth(ctx)
		if !health.IsHealthy {
			h.IsHealthy = false
			h.Info[StatusKey] = UnHealthyStatus
		} else {
			h.Info[StatusKey] = HealthyStatus
		}
		h.Info[checker.name] = health.Info
	}

	if h.IsHealthy {
		h.Info[StatusKey] = HealthyStatus
	} else {
		h.Info[StatusKey] = UnHealthyStatus
	}

	return h
}

func (c *Checker) Health(w http.ResponseWriter, req *http.Request) {
	health := c.IsHealthy()

	if !health.IsHealthy {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(health.Info)
}
