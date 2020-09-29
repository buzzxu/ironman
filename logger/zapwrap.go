package logger

import (
	"fmt"
	"go.uber.org/zap"
)

// CompatibleLogger is a logger which compatible to logrus/std log/prometheus.
// it implements Print() Println() Printf() Dbug() Debugln() Debugf() Info() Infoln() Infof() Warn() Warnln() Warnf()
// Error() Errorln() Errorf() Fatal() Fataln() Fatalf() Panic() Panicln() Panicf() With() WithField() WithFields()
type CompatibleLogger struct {
	_log *zap.Logger
}

func (l CompatibleLogger) Unwrap() *zap.Logger {
	return l._log
}

// Print logs a message at level Info on the compatibleLogger.
func (l CompatibleLogger) Print(args ...interface{}) {
	l._log.Info(fmt.Sprint(args...))
}

// Printf logs a message at level Info on the compatibleLogger.
func (l CompatibleLogger) Printf(format string, args ...interface{}) {
	l._log.Info(fmt.Sprintf(format, args...))
}

// Debug logs a message at level Debug on the compatibleLogger.
func (l CompatibleLogger) Debug(args ...interface{}) {
	l._log.Debug(fmt.Sprint(args...))
}

// Debugf logs a message at level Debug on the compatibleLogger.
func (l CompatibleLogger) Debugf(format string, args ...interface{}) {
	l._log.Debug(fmt.Sprintf(format, args...))
}

// Info logs a message at level Info on the compatibleLogger.
func (l CompatibleLogger) Info(args ...interface{}) {
	l._log.Info(fmt.Sprint(args...))
}

// Infof logs a message at level Info on the compatibleLogger.
func (l CompatibleLogger) Infof(format string, args ...interface{}) {
	l._log.Info(fmt.Sprintf(format, args...))
}

// Warn logs a message at level Warn on the compatibleLogger.
func (l CompatibleLogger) Warn(args ...interface{}) {
	l._log.Warn(fmt.Sprint(args...))
}

// Warnf logs a message at level Warn on the compatibleLogger.
func (l CompatibleLogger) Warnf(format string, args ...interface{}) {
	l._log.Warn(fmt.Sprintf(format, args...))
}

// Error logs a message at level Error on the compatibleLogger.
func (l CompatibleLogger) Error(args ...interface{}) {
	l._log.Error(fmt.Sprint(args...))
}

func (l CompatibleLogger) Errore(msg string, err error) {
	l._log.Error(msg, zap.Error(err))
}

// Errorf logs a message at level Error on the compatibleLogger.
func (l CompatibleLogger) Errorf(format string, args ...interface{}) {
	l._log.Error(fmt.Sprintf(format, args...))
}

// Fatal logs a message at level Fatal on the compatibleLogger.
func (l CompatibleLogger) Fatal(args ...interface{}) {
	l._log.Fatal(fmt.Sprint(args...))
}

// Fatalf logs a message at level Fatal on the compatibleLogger.
func (l CompatibleLogger) Fatalf(format string, args ...interface{}) {
	l._log.Fatal(fmt.Sprintf(format, args...))
}

// Panic logs a message at level Painc on the compatibleLogger.
func (l CompatibleLogger) Panic(args ...interface{}) {
	l._log.Panic(fmt.Sprint(args...))
}

// Panicf logs a message at level Painc on the compatibleLogger.
func (l CompatibleLogger) Panicf(format string, args ...interface{}) {
	l._log.Panic(fmt.Sprintf(format, args...))
}

// With return a logger with an extra field.
func (l *CompatibleLogger) With(key string, value interface{}) *CompatibleLogger {
	return &CompatibleLogger{l._log.With(zap.Any(key, value))}
}

// WithField return a logger with an extra field.
func (l *CompatibleLogger) WithField(key string, value interface{}) *CompatibleLogger {
	return &CompatibleLogger{l._log.With(zap.Any(key, value))}
}

// WithFields return a logger with extra fields.
func (l *CompatibleLogger) WithFields(fields map[string]interface{}) *CompatibleLogger {
	i := 0
	var clog *CompatibleLogger
	for k, v := range fields {
		if i == 0 {
			clog = l.WithField(k, v)
		} else {
			clog = clog.WithField(k, v)
		}
		i++
	}
	return clog
}
