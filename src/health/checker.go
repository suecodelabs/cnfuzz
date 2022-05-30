package health

import (
	"context"
	"encoding/json"
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
