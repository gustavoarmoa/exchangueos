// Package telemetry centralizes logging + OpenTelemetry tracing/metrics for ExchangeOS.
package telemetry

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger returns a structured zap logger. JSON in production, console in dev.
func NewLogger(env string) (*zap.Logger, error) {
	var cfg zap.Config
	if env == "production" || env == "staging" {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "ts"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.EncoderConfig.MessageKey = "msg"
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	return cfg.Build()
}
