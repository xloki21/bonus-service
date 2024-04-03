package log

import (
	"errors"
	"github.com/xloki21/bonus-service/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
	"time"
)

func NewZapLogger(level string, encoding string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	cfg.EncoderConfig.CallerKey = zapcore.OmitKey
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	cfg.OutputPaths = []string{"stdout"}
	cfg.Encoding = encoding
	parsedLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	cfg.Level = zap.NewAtomicLevelAt(parsedLevel)
	return cfg.Build()
}

var defaultLogger Logger
var once sync.Once

var TestLoggerConfig = config.LoggerConfig{
	Level:    "error",
	Encoding: "json",
}

func GetLogger() (Logger, error) {
	if defaultLogger == nil {
		return nil, errors.New("default logger is not initialized")
	}
	return defaultLogger, nil
}

func BuildLogger(cfg config.LoggerConfig) Logger {
	once.Do(func() {
		logger, err := NewZapLogger(cfg.Level, cfg.Encoding)
		if err != nil {
			panic(err)
		}
		defaultLogger = logger.Sugar()
	})
	return defaultLogger
}
