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
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, errors.New("")
	}
}

// InitLogger initializes the logger with the given parameters
func InitLogger(appName string, appID string, level string, dir string, color bool) error {
	logLevel, err := stringToLogLevel(level)
	if err != nil {
		return err
	}

	initLogger(appName, appID, logLevel, dir, color)
	return nil
}

func initLogger(appName string, appID string, level zapcore.Level, dir string, color bool) {
	// Construct production encoder config, set time format
	e := zap.NewDevelopmentEncoderConfig()
	e.EncodeTime = KoinosTimeEncoder
	if color {
		e.EncodeLevel = KoinosColorLevelEncoder
	} else {
		e.EncodeLevel = zapcore.LowercaseLevelEncoder
	}

	// Construct Console encoder for console output
	consoleEncoder := NewKoinosEncoder(e, appID)

	var coreFunc zap.Option

	if len(dir) > 0 {
		// Construct encoder for file output
		var fileEncoder zapcore.Encoder
		fe := zap.NewDevelopmentEncoderConfig()
		fe.EncodeTime = KoinosTimeEncoder
		fe.EncodeLevel = zapcore.LowercaseLevelEncoder
		fileEncoder = NewKoinosEncoder(fe, appID)

		// Construct lumberjack log roller
		lj := &lumberjack.Logger{
			Filename:   path.Join(dir, appName+".log"),
			MaxSize:    1,  // 1 Mb
			MaxBackups: 10, // 100 files
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
