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

// StringToLogLevel takes a string and returns a zap log level
func StringToLogLevel(level string) (zapcore.Level, error) {
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
func InitLogger(level zapcore.Level, jsonFileOutput bool, logFilename string, appID string) {
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

	// Construct logger
	logger, err := zap.NewProduction(coreFunc)
	if err != nil {
		panic(fmt.Sprintf("Error constructing logger: %v", err))
	}

	// Set global logger
	zap.ReplaceGlobals(logger)
}
