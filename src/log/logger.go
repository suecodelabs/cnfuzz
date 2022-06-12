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

package log

import (
	"fmt"
	"go.uber.org/zap/zaptest/observer"
	"log"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type ILogger interface {
	// Debug uses fmt.Sprint to construct and log a message.
	Debug(args ...any)
	// Info uses fmt.Sprint to construct and log a message.
	Info(args ...any)
	// Warn uses fmt.Sprint to construct and log a message.
	Warn(args ...any)
	// Error uses fmt.Sprint to construct and log a message.
	Error(args ...any)
	// DPanic uses fmt.Sprint to construct and log a message. In development, the
	// log then panics. (See DPanicLevel for details.)
	DPanic(args ...any)
	// Panic uses fmt.Sprint to construct and log a message, then panics.
	Panic(args ...any)
	// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
	Fatal(args ...any)
	// Debugf uses fmt.Sprintf to log a templated message.
	Debugf(template string, args ...any)
	// Infof uses fmt.Sprintf to log a templated message.
	Infof(template string, args ...any)
	// Warnf uses fmt.Sprintf to log a templated message.
	Warnf(template string, args ...any)
	// Errorf uses fmt.Sprintf to log a templated message.
	Errorf(template string, args ...any)
	// DPanicf uses fmt.Sprintf to log a templated message. In development, the
	// log then panics. (See DPanicLevel for details.)
	DPanicf(template string, args ...any)
	// Panicf uses fmt.Sprintf to log a templated message, then panics.
	Panicf(template string, args ...any)
	// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
	Fatalf(template string, args ...any)
	// Debugw logs a message with some additional context. The variadic key-value
	// pairs are treated as they are in With.
	//
	// When debug-level logging is disabled, this is much faster than
	//  s.With(keysAndValues).Debug(msg)
	Debugw(msg string, keysAndValues ...any)
	// Infow logs a message with some additional context. The variadic key-value
	// pairs are treated as they are in With.
	Infow(msg string, keysAndValues ...any)
	// Warnw logs a message with some additional context. The variadic key-value
	// pairs are treated as they are in With.
	Warnw(msg string, keysAndValues ...any)
	// Errorw logs a message with some additional context. The variadic key-value
	// pairs are treated as they are in With.
	Errorw(msg string, keysAndValues ...any)
	// DPanicw logs a message with some additional context. In development, the
	// log then panics. (See DPanicLevel for details.) The variadic key-value
	// pairs are treated as they are in With.
	DPanicw(msg string, keysAndValues ...any)
	// Panicw logs a message with some additional context, then panics. The
	// variadic key-value pairs are treated as they are in With.
	Panicw(msg string, keysAndValues ...any)
	// Fatalw logs a message with some additional context, then calls os.Exit. The
	// variadic key-value pairs are treated as they are in With.
	Fatalw(msg string, keysAndValues ...any)
}

var singleInstance ILogger

/* func CreateLogger(isDebug bool) ILogger {
	var zapErr error
	var log *zap.Logger
	if isDebug {
		log, zapErr = zap.NewDevelopment()
	} else {
		log, zapErr = zap.NewProduction()
	}

	if zapErr != nil {
		log.Fatalln(fmt.Errorf("error while trying to create a log instance: %w", zapErr).Error())
	}
	defer log.Sync()
	return log.Sugar()
} */

func SetupLogsCapture() *observer.ObservedLogs {
	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core)
	defer logger.Sync()
	singleInstance = logger.Sugar()
	return logs
}

// InitLogger initialize zap logger
func InitLogger(isDebug bool) {
	var zapErr error
	var logger *zap.Logger
	if isDebug {
		logger, zapErr = zap.NewDevelopment()
	} else {
		logger, zapErr = zap.NewProduction()
	}

	if zapErr != nil {
		log.Fatalln(fmt.Errorf("error while trying to create a log instance: %w", zapErr).Error())
	}
	defer logger.Sync()
	singleInstance = logger.Sugar()
}

// L get log instance
func L() ILogger {
	if singleInstance == nil {
		// FIXME: create a cleaner solution for testing for debug mode
		isDebug := viper.GetBool("debug")
		InitLogger(isDebug)
	}
	return singleInstance
}
