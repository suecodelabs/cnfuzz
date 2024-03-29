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
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"os"
)

var observedLogs *observer.ObservedLogs

// CreateDebugLogger creates a logger with settings optimized for debugging CnFuzz.
func CreateDebugLogger() Logger {
	l := zapr.NewLogger(createObservedZapLogger())
	e := CreateExiter(func(i int) {})
	return Logger{
		l,
		e,
	}
}

// GetObservedLogs returns ObservedLogs object, can be used in unit tests to see if something got logged.
func GetObservedLogs() *observer.ObservedLogs {
	return observedLogs
}

func (l Logger) GetExiter() *Exit {
	return l.exiter
}

func createObservedZapLogger() *zap.Logger {
	observedCore, logs := observer.New(zapcore.Level(PerformanceTestLevel)) // TODO implement variable info level
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		os.Stderr,
		zapcore.Level(PerformanceTestLevel),
	)
	zLogger := zap.New(zapcore.NewTee(observedCore, core))
	zap.ReplaceGlobals(zLogger)
	defer zLogger.Sync()
	observedLogs = logs
	return zLogger
}
