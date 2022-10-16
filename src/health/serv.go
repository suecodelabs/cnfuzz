// Copyright 2022 Sue B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package health

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Checker) health(w http.ResponseWriter, _ *http.Request) {
	health := c.IsHealthy()

	if !health.IsHealthy {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(health.Info)
}

func live(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprint(w, "ALIVE\n")
}

func (c *Checker) ready(w http.ResponseWriter, _ *http.Request) {
	health := c.IsHealthy()

	if !health.IsHealthy {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("NOT_READY\n"))
	} else {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("READY\n"))
	}
}

// Serv start http server that contains ready, live and health endpoints
// warning: this function is blocking
func Serv(hc Checker) {
	// Using this as a guideline
	// https://testfully.io/blog/api-health-check-monitoring/
	http.HandleFunc("/health", hc.health)
	http.HandleFunc("/health/live", live)
	http.HandleFunc("/health/ready", hc.ready)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		hc.l.FatalError(err, "failed to start webserver for health checks")
	}
}
