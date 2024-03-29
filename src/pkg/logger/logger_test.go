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

package logger

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFatalError(t *testing.T) {
	l := CreateDebugLogger()
	err := fmt.Errorf("test error")
	msg := "this is a test crash"
	l.FatalError(err, msg)
	assert.Equal(t, 1, l.GetExiter().status)
	logs := GetObservedLogs()
	theMsg := logs.FilterMessage(msg)
	assert.Equal(t, theMsg.All()[0].Message, msg)
}

func TestFatal(t *testing.T) {
	l := CreateDebugLogger()
	msg := "this is a test crash"
	l.Fatal(msg)
	assert.Equal(t, 1, l.GetExiter().status)
	logs := GetObservedLogs()
	theMsg := logs.FilterMessage(msg)
	assert.Equal(t, theMsg.All()[0].Message, msg)
}

func TestFatalKeysValues(t *testing.T) {
	l := CreateDebugLogger()
	msg := "this is a test crash"
	testKey := "key"
	testVal := "value"
	l.Fatal(msg, testKey, testVal)
	assert.Equal(t, 1, l.GetExiter().status)
	logs := GetObservedLogs()
	theMsg := logs.FilterMessage(msg)
	assert.Equal(t, theMsg.All()[0].Message, msg)
	assert.Equal(t, theMsg.All()[0].Context[0].Key, testKey)
	assert.Equal(t, theMsg.All()[0].Context[0].String, testVal)
}
