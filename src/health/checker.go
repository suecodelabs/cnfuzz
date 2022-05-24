package health

import (
	"encoding/json"
	"net/http"
)

type ICheck interface {
	CheckHealth() Health
}

type Check struct {
	name  string
	check ICheck
}

type Checker struct {
	checkers []Check
}

func (c Checker) IsHealthy() Health {
	h := NewHealth(true)

	for _, checker := range c.checkers {
		health := checker.check.CheckHealth()
		if !health.IsHealthy {
			h.IsHealthy = false
		}
		h.info[checker.name] = h.info
	}

	if h.IsHealthy {
		h.info["status"] = "pass"

	} else {
		h.info["status"] = "fail"
	}

	return h
}

func (c Checker) Health(w http.ResponseWriter, req *http.Request) {
	health := c.IsHealthy()

	if !health.IsHealthy {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(health.info)
}
