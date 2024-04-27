package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"go.uber.org/zap"
	"gopkg.in/errgo.v2/fmt/errors"
)

var (
	crowPackage string
	logger      *zap.Logger
	// Positions in the call stack when tracing to report the calling method
	minimumCallerDepth int
	callerInitOnce     sync.Once
)

const (
	maximumCallerDepth int = 25
	knownloggerFrames  int = 4
)

func initLogger() error {
	if config.Logger == nil {
		return nil
	}

	cfgs := []zap.Option{zap.AddCallerSkip(1)}
	if config.Logger.StackTrace {
		cfgs = append(cfgs, zap.AddStacktrace(zap.ErrorLevel))
	}
	if config.Logger.DevLogs {
		cfgs = append(cfgs, zap.IncreaseLevel(zap.InfoLevel))
	}

	var err error = nil
	logger, err = zap.NewProduction(cfgs...)
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Logger is constructed")
	return nil
}

func stopLogger() {
	if logger != nil {
		logger = nil
	}
}

func Data(v any) *zap.Logger {
	j, err := json.Marshal(v)
	if err != nil {
		return logger.With(zap.String(
			"data", fmt.Sprint(v),
		))
	}
	return logger.With(zap.String(
		"data", string(j),
	))
}

func Debug(args ...any) {
	if logger == nil {
		return
	}
	logger.Sugar().Debug(args...)
}

func Debugf(format string, args ...any) {
	if logger == nil {
		return
	}
	logger.Sugar().Debugf(format, args...)
}

func Info(args ...any) {
	if logger == nil {
		return
	}
	logger.Sugar().Info(args...)
}

func Infof(format string, args ...any) {
	if logger == nil {
		return
	}
	logger.Sugar().Infof(format, args...)
}

func Warn(args ...any) {
	if logger == nil {
		return
	}
	logger.Sugar().Warn(args...)
}

func Warnf(format string, args ...any) {
	if logger == nil {
		return
	}
	logger.Sugar().Warnf(format, args...)
}

func Error(args ...any) error {
	if logger == nil {
		return nil
	}
	logger.Sugar().Error(args...)
	return errors.New(fmt.Sprint(args...))
}

func Errorf(format string, args ...any) error {
	if logger == nil {
		return nil
	}
	logger.Sugar().Errorf(format, args...)
	return errors.New(fmt.Sprintf(format, args...))
}

func Fatal(args ...any) {
	if logger == nil {
		os.Exit(1)
	}
	logger.Sugar().Fatal(args...)
}

func Fatalf(format string, args ...any) {
	if logger == nil {
		os.Exit(1)
	}
	logger.Sugar().Fatalf(format, args...)
}

func Panic(args ...any) {
	if logger == nil {
		panic(args)
	}
	logger.Sugar().Panic(args...)
}

func Panicf(format string, args ...any) {
	if logger == nil {
		panic(args)
	}
	logger.Sugar().Panicf(format, args...)
}
