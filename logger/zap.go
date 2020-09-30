package logger

import (
	"github.com/buzzxu/ironman/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func zapLogger(level string) (logger *zap.Logger) {
	loggerConfig := newLoggerConfig(level, nil, nil)
	logger, err := loggerConfig.Build()
	if err != nil {
		panic(err)
	}
	return
}

func newZapLogger(level string, te zapcore.TimeEncoder, de zapcore.DurationEncoder) (logger *zap.Logger) {
	loggerConfig := newLoggerConfig(level, nil, nil)
	logger, err := loggerConfig.Build()
	if err != nil {
		panic(err)
	}
	return
}

func newLoggerConfig(level string, te zapcore.TimeEncoder, de zapcore.DurationEncoder) (loggerConfig zap.Config) {
	loggerConfig = zap.NewProductionConfig()
	if te == nil {
		loggerConfig.EncoderConfig.EncodeTime = timeEncoder
	} else {
		loggerConfig.EncoderConfig.EncodeTime = te
	}
	if de == nil {
		loggerConfig.EncoderConfig.EncodeDuration = milliSecondsDurationEncoder
	} else {
		loggerConfig.EncoderConfig.EncodeDuration = de
	}
	loggerConfig.EncoderConfig.EncodeCaller = callerEncoder
	switch level {
	case "debug":
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		break
	case "info":
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		break
	case "warn":
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
		break
	case "error":
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
		break
	default:
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	return
}

func newLogger(level string, json, console, caller bool, lumberjack *lumberjack.Logger) (logger *zap.Logger) {
	writeSyncer := getLogWriter(console, lumberjack)
	encoder := getEncoder(json, caller)
	var logLevel zapcore.Level
	switch strings.ToLower(level) {
	case "debug":
		logLevel = zap.DebugLevel
		break
	case "info":
		logLevel = zap.InfoLevel
		break
	case "warn":
		logLevel = zap.WarnLevel
		break
	case "error":
		logLevel = zap.ErrorLevel
		break
	default:
		logLevel = zap.DebugLevel
	}
	core := zapcore.NewCore(encoder, writeSyncer, logLevel)
	logger = zap.New(core, zap.AddCaller())
	return
}

func defaultLogger(json, console, caller bool) (logger *zap.Logger) {
	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		//不显示debug
		return lev < zap.ErrorLevel && lev > zap.DebugLevel
	})
	encoder := getEncoder(json, caller)
	lowWriteSyncer := getLogWriter(console, lumberJack(conf.ServerConf.Logger.Dir+string(filepath.Separator)+"info.log"))
	highWriteSyncer := getLogWriter(console, lumberJack(conf.ServerConf.Logger.Dir+string(filepath.Separator)+"error.log"))

	lowCore := zapcore.NewCore(encoder, lowWriteSyncer, lowPriority)
	highCore := zapcore.NewCore(encoder, highWriteSyncer, highPriority)
	logger = zap.New(zapcore.NewTee(highCore, lowCore), zap.AddCaller())
	return
}

func getEncoder(json, caller bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.StacktraceKey = "stacktrace"
	encoderConfig.CallerKey = "line"
	encoderConfig.EncodeTime = timeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	if caller {
		encoderConfig.EncodeCaller = callerEncoder
	}
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeName = zapcore.FullNameEncoder
	if json {
		return zapcore.NewJSONEncoder(encoderConfig)
	} else {
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
}
func getLogWriter(console bool, hook *lumberjack.Logger) zapcore.WriteSyncer {
	if console {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(hook))
	}
	return zapcore.AddSync(hook)
}

// callerEncoder will add caller to log. format is "filename:lineNum:funcName", e.g:"zaplog/zaplog_test.go:15:zaplog.TestNewLogger"
func callerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(strings.Join([]string{caller.TrimmedPath(), runtime.FuncForPC(caller.PC).Name()}, ":"))
}

// timeEncoder specifics the time format
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// milliSecondsDurationEncoder serializes a time.Duration to a floating-point number of micro seconds elapsed.
func milliSecondsDurationEncoder(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendFloat64(float64(d) / float64(time.Millisecond))
}
