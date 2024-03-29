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
	"github.com/suecodelabs/cnfuzz/src/pkg/logger"
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
	l        logger.Logger
	checkers []Check
}

func NewChecker(l logger.Logger) Checker {
	return Checker{
		l:        l,
		checkers: make([]Check, 0),
	}
}

func (c *Checker) RegisterCheck(name string, check ICheck) {
	if check == nil {
		c.l.V(logger.ImportantLevel).Info("failed to register health check, because it doesn't contain a check function", "checkName", name)
		return
	}
	newCheck := Check{name: name, check: check}
	c.checkers = append(c.checkers, newCheck)
}

func (c *Checker) IsHealthy() Health {
	h := NewHealth(true)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
	defer cancel()

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
