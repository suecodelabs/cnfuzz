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
	"os"
)

const (
	ImportantLevel       = int(zapcore.ErrorLevel)
	InfoLevel            = int(zapcore.InfoLevel)
	DebugLevel           = int(zapcore.DebugLevel)
	PerformanceTestLevel = -128
)

type Logger struct {
	logr.Logger
	exiter *Exit
}

// CreateLogger creates a new logger instance
// isDebug: if enabled, the logger prints debug logs, otherwise it prints info level and above
func CreateLogger(isDebug bool, logLevel int) Logger {
	var logger logr.Logger
	logger = zapr.NewLogger(createZapLogger(isDebug, logLevel))
	return Logger{
		logger,
		nil,
	}
}

// SetExiter useful for debugging
// when the exit is set, the code won't call os.Exit($code) but will set the status on the Exit struct
// this is useful when you want to test if a function crashes with os.Exit(1) for example.
func (l Logger) SetExiter(exit *Exit) {
	l.exiter = exit
}

// Fatal logs message with important level and exits with code 1
func (l Logger) Fatal(msg string, keysAndValues ...interface{}) {
	if keysAndValues == nil {
		l.V(ImportantLevel).Info(msg)
	} else {
		l.V(ImportantLevel).Info(msg, keysAndValues...)
	}
	l.exiter.Exit(1)
}

// FatalError logs message and error with important level and exits with code 1
func (l Logger) FatalError(err error, msg string, keysAndValues ...interface{}) {
	if keysAndValues == nil {
		l.V(ImportantLevel).Error(err, msg)
	} else {
		l.V(ImportantLevel).Error(err, msg, keysAndValues...)
	}
	l.exiter.Exit(1)
}

func createZapLogger(isDebug bool, logLevel int) *zap.Logger {
	var zLogger *zap.Logger
	var err error
	if isDebug {
		core := zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			os.Stderr,
			zapcore.Level(logLevel),
		)
		zLogger = zap.New(core)
	} else {
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			os.Stderr,
			zapcore.Level(logLevel),
		)
		zLogger = zap.New(core)
	}
	if err != nil {
		log.Fatalf("failed to create logger:\n%s", err)
	}

	defer zLogger.Sync()

	return zLogger
}
