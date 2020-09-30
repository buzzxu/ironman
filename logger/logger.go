package logger

import (
	"github.com/buzzxu/ironman/conf"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"path/filepath"
	"strings"
)

var logs map[string]*CompatibleLogger
var logger *CompatibleLogger

func InitLogger() {
	//设置默认logger
	logger = &CompatibleLogger{defaultLogger(conf.ServerConf.Logger.Json, conf.ServerConf.Logger.Console, true).WithOptions(zap.AddCallerSkip(1))}
	logs = make(map[string]*CompatibleLogger)
	for name, logger := range conf.ServerConf.Logger.Loggers {
		if logger.Level == "" {
			logger.Level = "info"
		}
		if logger.File == "" {
			logger.File = conf.ServerConf.Logger.Dir + string(filepath.Separator) + name + "." + logger.Level + ".log"
		}
		if !strings.HasPrefix(logger.File, "./") && !strings.HasPrefix(logger.File, "/") {
			logger.File = conf.ServerConf.Logger.Dir + string(filepath.Separator) + logger.File
		}
		logs[name] = newCompatibleLogger(name, logger)
	}
}

//New
func New(name, level string, json, console, caller bool) (*CompatibleLogger, error) {
	return NewLogger(name, name+".log", level, json, console, caller)
}

//NewLogger
func NewLogger(name, file, level string, json, console, caller bool) (*CompatibleLogger, error) {
	if name == "" {
		return nil, errors.New("name is empty")
	}
	if level == "" {
		level = "info"
	}
	if file == "" {
		file = conf.ServerConf.Logger.Dir + string(filepath.Separator) + name + "." + level + ".log"
	}
	if !strings.HasPrefix(file, "./") && !strings.HasPrefix(file, "/") {
		file = conf.ServerConf.Logger.Dir + string(filepath.Separator) + file
	}
	return &CompatibleLogger{newLogger(level, json, console, caller, lumberJack(file)).WithOptions(zap.AddCallerSkip(1))}, nil
}

func Logger(name string) (log *CompatibleLogger) {
	var ok bool
	log, ok = logs[name]
	//如果不存在返回默认log
	if !ok {
		log = logger
	}
	return log
}

func Of(name string) (log *CompatibleLogger) {
	return Logger(name)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}
func Warn(args ...interface{}) {
	logger.Warn(args...)
}
func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}
func Error(args ...interface{}) {
	logger.Error(args...)
}
func Errore(msg string, err error) {
	logger.Errore(msg, err)
}
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}
func Debug(args ...interface{}) {
	logger.Debug(args...)
}
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}
func Fatale(msg string, err error) {
	logger.Fatale(msg, err)
}
func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

func X() *CompatibleLogger {
	return logger
}

func lumberJack(file string) (hook *lumberjack.Logger) {
	hook = &lumberjack.Logger{
		Filename:   file,                              // 日志文件位置
		MaxSize:    conf.ServerConf.Logger.MaxSize,    // 日志文件最大大小(MB)
		MaxBackups: conf.ServerConf.Logger.MaxBackups, // 保留旧文件最大数量
		Compress:   conf.ServerConf.Logger.Compress,   // 是否压缩旧文件
	}
	if conf.ServerConf.Logger.MaxAge > 0 {
		hook.MaxAge = conf.ServerConf.Logger.MaxAge // days 保留旧文件最长天数
	}
	return
}

//newCompatibleLogger
func newCompatibleLogger(name string, logConf *conf.LogConf) *CompatibleLogger {
	var (
		console = conf.ServerConf.Logger.Console
		json    = conf.ServerConf.Logger.Json
	)
	if console != logConf.Console {
		console = logConf.Console
	}

	if json != logConf.Json {
		json = logConf.Json
	}
	l, err := NewLogger(name, logConf.File, logConf.Level, json, console, true)
	if err != nil {
		log.Fatalf("logger[%s] 创建失败: %v", name, err)
	}
	return l
}
