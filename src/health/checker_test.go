package health

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewChecker(t *testing.T) {
	c := NewChecker()
	assert.NotNil(t, c.checkers)
}

func TestRegisterHealthyCheck(t *testing.T) {
	c := NewChecker()
	check := fakeHealthyCheck{}
	checkName := "test-checker"
	c.RegisterCheck(checkName, check)
	status := c.IsHealthy()
	assert.True(t, status.IsHealthy)
	assert.Equal(t, HealthyStatus, status.Info[StatusKey])
	var testCheckMap map[string]any
	testCheckMap = status.Info[checkName].(map[string]any)
	assert.Equal(t, HealthyStatus, testCheckMap[StatusKey])
}

type fakeHealthyCheck struct {
}

func (fakeHealthyCheck) CheckHealth(context context.Context) Health {
	h := NewHealth(true)
	h.Info[StatusKey] = HealthyStatus
	return h
}

func TestRegisterUnhealthyCheck(t *testing.T) {
	c := NewChecker()
	check := fakeUnhealthyCheck{}
	checkName := "test-checker"
	c.RegisterCheck(checkName, check)
	status := c.IsHealthy()
	assert.False(t, status.IsHealthy)
	assert.Equal(t, UnHealthyStatus, status.Info[StatusKey])
	var testCheckMap map[string]any
	testCheckMap = status.Info[checkName].(map[string]any)
	assert.Equal(t, UnHealthyStatus, testCheckMap[StatusKey])
}

type fakeUnhealthyCheck struct {
}

func (fakeUnhealthyCheck) CheckHealth(context context.Context) Health {
	h := NewHealth(false)
	h.Info[StatusKey] = UnHealthyStatus
	return h
}

func TestHttpHealthCheck(t *testing.T) {
	// TODO create more test cases
	c := NewChecker()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	c.Health(w, req)
	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	data, err := ioutil.ReadAll(res.Body)
	if assert.NoError(t, err, "http health request shouldn't error out") {
		var payload map[string]interface{}
		err = json.Unmarshal(data, &payload)
		if assert.NoError(t, err, "returned json is invalid") {
			assert.Equal(t, HealthyStatus, payload[StatusKey])
		}
	}
}
