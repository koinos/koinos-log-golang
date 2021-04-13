package log

import "go.uber.org/zap"

func Debugf(template string, args ...interface{}) {
	zap.S().Debugf(template, args...)
}

func Debug(msg string) {
	zap.L().Debug(msg)
}

func Infof(template string, args ...interface{}) {
	zap.S().Infof(template, args...)
}

func Info(msg string) {
	zap.L().Info(msg)
}

func Warnf(template string, args ...interface{}) {
	zap.S().Warnf(template, args...)
}

func Warn(msg string) {
	zap.L().Warn(msg)
}

func Errorf(template string, args ...interface{}) {
	zap.S().Errorf(template, args...)
}

func Error(msg string) {
	zap.L().Error(msg)
}
