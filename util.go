package log

import (
	"errors"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Log consts
const (
	maxSize         = 128
	maxBackups      = 32
	maxAge          = 64
	compressBackups = true
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

func InitLogger(level string, jsonFileOutput bool, logFilename string, appID string) error {
	l, err := stringToLogLevel(level)
	if err != nil {
		return err
	}

	initLogger(l, jsonFileOutput, logFilename, appID)
	return nil
}

// InitLogger initializes the logger with the given parameters
func initLogger(level zapcore.Level, jsonFileOutput bool, logFilename string, appID string) {
	// Construct production encoder config, set time format
	e := zap.NewDevelopmentEncoderConfig()
	e.EncodeTime = KoinosTimeEncoder
	e.EncodeLevel = KoinosColorLevelEncoder

	// Construct encoder for file output
	var fileEncoder zapcore.Encoder
	if jsonFileOutput { // Json encoder
		fileEncoder = zapcore.NewJSONEncoder(e)
	} else { // Console encoder, minus log-level coloration
		fe := zap.NewDevelopmentEncoderConfig()
		fe.EncodeTime = KoinosTimeEncoder
		fe.EncodeLevel = zapcore.LowercaseLevelEncoder
		fileEncoder = NewKoinosEncoder(fe, appID)
	}

	// Construct Console encoder for console output
	consoleEncoder := NewKoinosEncoder(e, appID)

	// Construct lumberjack log roller
	lj := &lumberjack.Logger{
		Filename:   logFilename,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compressBackups,
	}

	// Construct core
	coreFunc := zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return zapcore.NewTee(
			zapcore.NewCore(fileEncoder, zapcore.AddSync(lj), level),
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
		)
	})

	// Construct logger. Add caller skip for correct line numbers (since this library wraps zap calls)
	logger, err := zap.NewProduction(coreFunc, zap.AddCallerSkip(1))
	if err != nil {
		panic(fmt.Sprintf("Error constructing logger: %v", err))
	}

	// Set global logger
	zap.ReplaceGlobals(logger)
}
