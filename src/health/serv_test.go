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
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/suecodelabs/cnfuzz/src/logger"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpHealth(t *testing.T) {
	cases := []struct {
		check              ICheck
		checkName          string
		expectedStatus     string
		expectedStatusCode int
	}{
		{FakeUnhealthyCheck{}, "unhealthy-test-check", UnHealthyStatus, http.StatusServiceUnavailable},
		{FakeHealthyCheck{}, "healthy-test-check", HealthyStatus, http.StatusOK},
		{nil, "nil-test-check", HealthyStatus, http.StatusOK},
	}

	for _, c := range cases {
		l := logger.CreateDebugLogger()
		checker := NewChecker(l)

		checker.RegisterCheck(c.checkName, c.check)

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()
		checker.health(w, req)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, c.expectedStatusCode, res.StatusCode, "http request didn't return the expected status code")

		data, err := io.ReadAll(res.Body)
		if assert.NoError(t, err, "reading http health request body shouldn't error out") {
			var respObj map[string]interface{}
			err = json.Unmarshal(data, &respObj)
			if assert.NoError(t, err, "returned json is invalid") {
				assert.Equal(t, c.expectedStatus, respObj[StatusKey])
			}
		}
	}
}

func TestHttpLive(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	w := httptest.NewRecorder()
	live(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if assert.NoError(t, err, "reading http live request body shouldn't error out") {
		content := string(data[:])
		assert.Equal(t, "ALIVE\n", content)
	}
}

func TestHttpReady(t *testing.T) {
	cases := []struct {
		check              ICheck
		checkName          string
		expectedStatus     string
		expectedStatusCode int
	}{
		{FakeUnhealthyCheck{}, "unhealthy-test-check", "NOT_READY\n", http.StatusServiceUnavailable},
		{FakeHealthyCheck{}, "healthy-test-check", "READY\n", http.StatusOK},
		{nil, "nil-test-check", "READY\n", http.StatusOK},
	}

	l := logger.CreateDebugLogger()
	for _, c := range cases {
		checker := NewChecker(l)
		checker.RegisterCheck(c.checkName, c.check)

		req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
		w := httptest.NewRecorder()
		checker.ready(w, req)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, c.expectedStatusCode, res.StatusCode, "http request didn't return the expected status code")

		data, err := io.ReadAll(res.Body)
		if assert.NoError(t, err, "reading http ready request body shouldn't error out") {
			content := string(data[:])
			assert.Equal(t, c.expectedStatus, content)
		}
	}

}
