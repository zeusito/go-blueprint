package config

import (
	"os"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
)

func NewLogger() *zap.Logger {
	encoderConfig := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encoderConfig, os.Stdout, zap.DebugLevel)
	logger := zap.New(core, zap.AddCaller())

	return logger
}
