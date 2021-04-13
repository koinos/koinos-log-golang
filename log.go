package log

import "go.uber.org/zap"

// Debugf logs to the debug level with a format string and associated parameters
func Debugf(template string, args ...interface{}) {
	zap.S().Debugf(template, args...)
}

// Debug logs to the debug level
func Debug(msg string) {
	zap.L().Debug(msg)
}

// Infof logs to the info level with a format string and associated parameters
func Infof(template string, args ...interface{}) {
	zap.S().Infof(template, args...)
}

// Info logs to the info level
func Info(msg string) {
	zap.L().Info(msg)
}

// Warnf logs to the warn level with a format string and associated parameters
func Warnf(template string, args ...interface{}) {
	zap.S().Warnf(template, args...)
}

// Warn logs to the warn level
func Warn(msg string) {
	zap.L().Warn(msg)
}

// Errorf logs to the error level with a format string and associated parameters
func Errorf(template string, args ...interface{}) {
	zap.S().Errorf(template, args...)
}

// Error logs to the error level
func Error(msg string) {
	zap.L().Error(msg)
}
