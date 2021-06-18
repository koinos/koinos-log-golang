package log

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
)

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

func FormatError(msg string, vars ...interface{}) string {
	if (len(vars) & 1) != 0 {
		Errorm("Message did not have even number of arguments", "msg", msg)
	}
	result := msg
	for i := 0; i < len(vars); i += 2 {
		name, ok := vars[i].(string)
		if !ok {
			Errorm("Message did not specify name as string", "msg", msg)
			break
		}
		result += fmt.Sprintf(" - %s is %v", name, vars[i+1])
	}
	return result
}

// Debugm logs to the debug level with consistent logging
func Debugm(msg string, args ...interface{}) {
	Debug(FormatError(msg, args))
}

// Infom logs to the info level with consistent logging
func Infom(msg string, args ...interface{}) {
	Info(FormatError(msg, args))
}

// Warnm logs to the warn level with consistent logging
func Warnm(msg string, args ...interface{}) {
	Warn(FormatError(msg, args))
}

// Errorm logs to the error level with consistent logging
func Errorm(msg string, args ...interface{}) {
	Error(FormatError(msg, args))
}

func NewErrorm(msg string, args ...interface{}) error {
	return errors.New(FormatError(msg, args))
}
