package health

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/suecodelabs/cnfuzz/src/log"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeHealthyCheck struct {
}

func (fakeHealthyCheck) CheckHealth(context context.Context) Health {
	h := NewHealth(true)
	h.Info[StatusKey] = HealthyStatus
	return h
}

type fakeUnhealthyCheck struct {
}

func (fakeUnhealthyCheck) CheckHealth(context context.Context) Health {
	h := NewHealth(false)
	h.Info[StatusKey] = UnHealthyStatus
	return h
}

func TestNewChecker(t *testing.T) {
	c := NewChecker()
	assert.NotNil(t, c.checkers)
}

func TestRegisterCheck(t *testing.T) {
	cases := []struct {
		check          ICheck
		checkName      string
		isHealthy      bool
		expectedStatus string
	}{
		{fakeUnhealthyCheck{}, "unhealthy-test-check", false, UnHealthyStatus},
		{fakeHealthyCheck{}, "healthy-test-check", true, HealthyStatus},
		{nil, "nil-test-check", true, HealthyStatus},
	}

	for _, c := range cases {
		checker := NewChecker()

		var status Health
		if c.check == nil {
			// test if the logger logs the failed health check registration

			logs := log.SetupLogsCapture()

			checker.RegisterCheck(c.checkName, c.check)
			status = checker.IsHealthy()

			entry := logs.All()[0]
			assert.Equal(t, zap.WarnLevel, entry.Level, "expected a message with the 'warn' log level got %s", entry.Level)
			assert.Equal(t, fmt.Sprintf("failed to register %s health check, because it doesn't contain a check function", c.checkName), entry.Message, "got an unexpected log message while registering a nil health check")
		} else {
			// in this scenario we don't expect any log messages, so just register the health check as usual
			checker.RegisterCheck(c.checkName, c.check)
			status = checker.IsHealthy()
		}

		assert.Equal(t, c.isHealthy, status.IsHealthy, "expected different health state")
		assert.Equal(t, c.expectedStatus, status.Info[StatusKey], "expected different health status key")
		if c.check != nil {
			var testCheckMap map[string]any
			testCheckMap = status.Info[c.checkName].(map[string]any)
			assert.Equal(t, c.expectedStatus, testCheckMap[StatusKey])
		}
	}

}

func TestHttpHealthCheck(t *testing.T) {
	cases := []struct {
		check              ICheck
		checkName          string
		expectedStatus     string
		expectedStatusCode int
	}{
		{fakeUnhealthyCheck{}, "unhealthy-test-check", UnHealthyStatus, http.StatusServiceUnavailable},
		{fakeHealthyCheck{}, "healthy-test-check", HealthyStatus, http.StatusOK},
		{nil, "nil-test-check", HealthyStatus, http.StatusOK},
	}

	for _, c := range cases {
		// TODO create more test cases
		checker := NewChecker()

		checker.RegisterCheck(c.checkName, c.check)

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()
		checker.Health(w, req)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, c.expectedStatusCode, res.StatusCode, "http request didn't return the expected status code")

		data, err := ioutil.ReadAll(res.Body)
		if assert.NoError(t, err, "http health request shouldn't error out") {
			var respObj map[string]interface{}
			err = json.Unmarshal(data, &respObj)
			if assert.NoError(t, err, "returned json is invalid") {
				assert.Equal(t, c.expectedStatus, respObj[StatusKey])
			}
		}

	}

}
