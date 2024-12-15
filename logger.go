package logger

import (
	"github.com/null-bd/logger/internal/core"
	"github.com/null-bd/logger/types"
)

// Expose types from types package
type (
	Level  = types.Level
	Fields = types.Fields
	Logger = types.Logger
	Config = types.Config
)

// Expose constants
const (
	DebugLevel = types.DebugLevel
	InfoLevel  = types.InfoLevel
	WarnLevel  = types.WarnLevel
	ErrorLevel = types.ErrorLevel
	FatalLevel = types.FatalLevel
)

// New creates a new logger instance with the provided configuration
func New(cfg *Config) (Logger, error) {
	return core.NewLogger(cfg)
}
