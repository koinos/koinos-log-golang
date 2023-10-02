package log

import (
	"errors"
	"fmt"
	"os"
	"path"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func stringToLogLevel(level string) (zapcore.Level, error) {
	switch level {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warning":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, errors.New("")
	}
}

// InitLogger initializes the logger with the given parameters
func InitLogger(appName string, instanceID string, level string, dir string, color bool, datetime bool) error {
	logLevel, err := stringToLogLevel(level)
	if err != nil {
		return err
	}

	initLogger(appName, instanceID, logLevel, dir, color, datetime)
	return nil
}

func initLogger(appName string, instanceID string, level zapcore.Level, dir string, color bool, datetime bool) {
	appID := fmt.Sprintf("%s.%s", appName, instanceID)

	// Construct production encoder config, set time format
	e := zap.NewDevelopmentEncoderConfig()
	if datetime {
		e.EncodeTime = KoinosTimeEncoder
	}

	if color {
		e.EncodeLevel = KoinosColorLevelEncoder
	} else {
		e.EncodeLevel = KoinosLevelEncoder
	}

	// Construct Console encoder for console output
	consoleEncoder := NewKoinosEncoder(e, appID)

	var coreFunc zap.Option

	if len(dir) > 0 {
		// Construct encoder for file output
		var fileEncoder zapcore.Encoder
		fe := zap.NewDevelopmentEncoderConfig()
		fe.EncodeTime = KoinosTimeEncoder
		fe.EncodeLevel = KoinosLevelEncoder
		fileEncoder = NewKoinosEncoder(fe, appID)

		// Construct lumberjack log roller
		lj := &lumberjack.Logger{
			Filename:   path.Join(dir, appName+".log"),
			MaxSize:    1,   // 1 Mb
			MaxBackups: 100, // 100 files
		}

		// Construct core
		coreFunc = zap.WrapCore(func(zapcore.Core) zapcore.Core {
			return zapcore.NewTee(
				zapcore.NewCore(fileEncoder, zapcore.AddSync(lj), level),
				zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
			)
		})
	} else {
		coreFunc = zap.WrapCore(func(zapcore.Core) zapcore.Core {
			return zapcore.NewTee(
				zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
			)
		})
	}

	// Construct logger. Add caller skip for correct line numbers (since this library wraps zap calls)
	logger, err := zap.NewProduction(coreFunc, zap.AddCallerSkip(1))
	if err != nil {
		panic(fmt.Sprintf("Error constructing logger: %v", err))
	}

	// Set global logger
	zap.ReplaceGlobals(logger)
}
