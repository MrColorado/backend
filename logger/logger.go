package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"log/syslog"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"gopkg.in/errgo.v2/fmt/errors"
)

var (
	crowPackage string
	logger      *logrus.Entry
	// Positions in the call stack when tracing to report the calling method
	minimumCallerDepth int
	callerInitOnce     sync.Once
)

const (
	maximumCallerDepth int = 25
	knownCrowFrames    int = 4
)

func initLogger() error {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetLevel(logrus.DebugLevel)

	if !config.Logger.DevLogs {
		if config.Logger.SysLogs {
			hook, err := syslogger.NewSyslogHook("", "", syslog.LOG_LOCAL1, config.AppName)
			if err != nil {
				return Error("Unable to connect to local syslog daemon")
			}
			logrus.AddHook(hook)
			logrus.SetOutput(io.Discard)
		}
		logrus.SetFormatter(&logrus.JSONFormatter{
			DisableTimestamp: true,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyMsg: "message",
			},
		})
		logrus.SetLevel(logrus.InfoLevel)
	}

	if useSentry {
		if config.Logger.SentryDsn == "" {
			useSentry = false
			return Error("crow: Cannot use Sentry as the SentryDsn configuration field is empty")
		}
		hook, err := logrus_sentry.NewSentryHook(config.Logger.SentryDsn, []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
		})
		if err != nil {
			return err
		}
		logrus.AddHook(hook)
	}

	return nil
}

func stopLogger() {
	if logger != nil {
		logger = nil
	}
}

// getPackageName reduces a fully qualified function name to the package name
// There really ought to be to be a better way...
func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}

// getCaller retrieves the name of the first non-crow calling function
func getCrowCaller() *runtime.Frame {
	// cache this package's fully-qualified name
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, maximumCallerDepth)
		_ = runtime.Callers(0, pcs)

		// dynamic get the package name and the minimum caller depth
		for i := 0; i < maximumCallerDepth; i++ {
			funcName := runtime.FuncForPC(pcs[i]).Name()
			if strings.Contains(funcName, "getCrowCaller") {
				crowPackage = getPackageName(funcName)
				break
			}
		}

		minimumCallerDepth = knownCrowFrames
	})

	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		// If the caller isn't part of this package, we're done
		if pkg != crowPackage {
			return &f //nolint:scopelint
		}
	}

	// if we got here, we failed to find the caller's context
	return nil
}

func getLogger() *logrus.Entry {
	if config.Logger.CallerLogs {
		if crowFrame := getCrowCaller(); crowFrame != nil {
			_, filename := path.Split(crowFrame.File)
			return logger.WithFields(logrus.Fields{
				"line": strconv.Itoa(crowFrame.Line),
				"file": filename,
			})
		}
	}
	return logger
}

func Data(v any) *logrus.Entry {
	j, err := json.Marshal(v)
	if err != nil {
		return getLogger().WithFields(logrus.Fields{
			"data": fmt.Sprint(v),
		})
	}
	return getLogger().WithFields(logrus.Fields{
		"data": string(j),
	})
}

func Debug(args ...any) {
	if logger == nil {
		return
	}
	getLogger().Debug(args...)
}

func Debugf(format string, args ...any) {
	if logger == nil {
		return
	}
	getLogger().Debugf(format, args...)
}

func Info(args ...any) {
	if logger == nil {
		return
	}
	getLogger().Info(args...)
}

func Infof(format string, args ...any) {
	if logger == nil {
		return
	}
	getLogger().Infof(format, args...)
}

func Warn(args ...any) {
	if logger == nil {
		return
	}
	getLogger().Warn(args...)
}

func Warnf(format string, args ...any) {
	if logger == nil {
		return
	}
	getLogger().Warnf(format, args...)
}

func Error(args ...any) error {
	if logger == nil {
		return nil
	}
	getLogger().Error(args...)
	return errors.New(fmt.Sprint(args...))
}

func Errorf(format string, args ...any) error {
	if logger == nil {
		return nil
	}
	getLogger().Errorf(format, args...)
	return errors.New(fmt.Sprintf(format, args...))
}

func Fatal(args ...any) {
	if logger == nil {
		os.Exit(1)
	}
	getLogger().Fatal(args...)
}

func Fatalf(format string, args ...any) {
	if logger == nil {
		os.Exit(1)
	}
	getLogger().Fatalf(format, args...)
}

func Panic(args ...any) {
	if logger == nil {
		panic(args)
	}
	getLogger().Panic(args...)
}

func Panicf(format string, args ...any) {
	if logger == nil {
		panic(args)
	}
	getLogger().Panicf(format, args...)
}
