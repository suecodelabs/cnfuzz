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
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

const (
	ImportantLevel       = int(zapcore.FatalLevel)
	InfoLevel            = int(zapcore.InfoLevel)
	DebugLevel           = int(zapcore.DebugLevel)
	PerformanceTestLevel = -128
)

// CreateLogger creates a new logger instance
// isDebug: if enabled, the logger prints debug logs, otherwise it prints info level and above
func CreateLogger(isDebug bool) logr.Logger {
	var logger logr.Logger
	logger = zapr.NewLogger(createZapLogger(isDebug))
	return logger
}

func createZapLogger(isDebug bool) *zap.Logger {
	var zLogger *zap.Logger
	var err error
	if isDebug {
		/* zc := zap.NewDevelopmentConfig()
		zc.Level = zap.NewAtomicLevelAt(DebugLevel)
		zLogger, err = zc.Build() */
		zLogger, err = zap.NewDevelopment() // defaults to debug level
	} else {
		zLogger, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatalf("failed to create logger:\n%s", err)
	}

	defer zLogger.Sync()

	return zLogger
}
