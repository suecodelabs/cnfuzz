package log

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"github.com/suecodelabs/cnfuzz/src/cmd"
	"go.uber.org/zap"
)

type ILogger interface {
	// Debug uses fmt.Sprint to construct and log a message.
	Debug(args ...interface{})
	// Info uses fmt.Sprint to construct and log a message.
	Info(args ...interface{})
	// Warn uses fmt.Sprint to construct and log a message.
	Warn(args ...interface{})
	// Error uses fmt.Sprint to construct and log a message.
	Error(args ...interface{})
	// DPanic uses fmt.Sprint to construct and log a message. In development, the
	// log then panics. (See DPanicLevel for details.)
	DPanic(args ...interface{})
	// Panic uses fmt.Sprint to construct and log a message, then panics.
	Panic(args ...interface{})
	// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
	Fatal(args ...interface{})
	// Debugf uses fmt.Sprintf to log a templated message.
	Debugf(template string, args ...interface{})
	// Infof uses fmt.Sprintf to log a templated message.
	Infof(template string, args ...interface{})
	// Warnf uses fmt.Sprintf to log a templated message.
	Warnf(template string, args ...interface{})
	// Errorf uses fmt.Sprintf to log a templated message.
	Errorf(template string, args ...interface{})
	// DPanicf uses fmt.Sprintf to log a templated message. In development, the
	// log then panics. (See DPanicLevel for details.)
	DPanicf(template string, args ...interface{})
	// Panicf uses fmt.Sprintf to log a templated message, then panics.
	Panicf(template string, args ...interface{})
	// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
	Fatalf(template string, args ...interface{})
	// Debugw logs a message with some additional context. The variadic key-value
	// pairs are treated as they are in With.
	//
	// When debug-level logging is disabled, this is much faster than
	//  s.With(keysAndValues).Debug(msg)
	Debugw(msg string, keysAndValues ...interface{})
	// Infow logs a message with some additional context. The variadic key-value
	// pairs are treated as they are in With.
	Infow(msg string, keysAndValues ...interface{})
	// Warnw logs a message with some additional context. The variadic key-value
	// pairs are treated as they are in With.
	Warnw(msg string, keysAndValues ...interface{})
	// Errorw logs a message with some additional context. The variadic key-value
	// pairs are treated as they are in With.
	Errorw(msg string, keysAndValues ...interface{})
	// DPanicw logs a message with some additional context. In development, the
	// log then panics. (See DPanicLevel for details.) The variadic key-value
	// pairs are treated as they are in With.
	DPanicw(msg string, keysAndValues ...interface{})
	// Panicw logs a message with some additional context, then panics. The
	// variadic key-value pairs are treated as they are in With.
	Panicw(msg string, keysAndValues ...interface{})
	// Fatalw logs a message with some additional context, then calls os.Exit. The
	// variadic key-value pairs are treated as they are in With.
	Fatalw(msg string, keysAndValues ...interface{})
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
		isDebug := viper.GetBool(cmd.IsDebug)
		InitLogger(isDebug)
	}
	return singleInstance
}
