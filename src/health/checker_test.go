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
	"github.com/stretchr/testify/assert"
	"github.com/suecodelabs/cnfuzz/src/logger"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestNewChecker(t *testing.T) {
	l := logger.CreateDebugLogger()
	c := NewChecker(l)
	assert.NotNil(t, c.checkers)
}

func TestRegisterCheck(t *testing.T) {
	cases := []struct {
		check          ICheck
		checkName      string
		isHealthy      bool
		expectedStatus string
	}{
		{FakeUnhealthyCheck{}, "unhealthy-test-check", false, UnHealthyStatus},
		{FakeHealthyCheck{}, "healthy-test-check", true, HealthyStatus},
		{nil, "nil-test-check", true, HealthyStatus},
	}

	for _, c := range cases {
		l := logger.CreateDebugLogger()
		checker := NewChecker(l)

		var status Health
		if c.check == nil {
			// test if the logger logs the failed health check registration

			checker.RegisterCheck(c.checkName, c.check)
			status = checker.IsHealthy()

			logs := logger.GetObservedLogs().All()
			entry := logs[0]
			assert.Equal(t, zapcore.Level(-logger.ImportantLevel), entry.Level, "expected a message with the 'important' log level got %s", entry.Level)
			assert.Equal(t, "failed to register health check, because it doesn't contain a check function", entry.Message, "got an unexpected log message while registering a nil health check")
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
