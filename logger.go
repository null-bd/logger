package logger

import (
	"github.com/null-bd/logger/internal/core"
	"github.com/null-bd/logger/types"
)

// Expose types from types package
type (
	Fields = types.Fields
	Level  = types.Level
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

// Expose trace functions
var (
	SetTraceFields   = core.SetTraceFields
	GetTraceFields   = core.GetTraceFields
	ClearTraceFields = core.ClearTraceFields
)

// New creates a new logger instance with the provided configuration
func New(cfg *Config) (Logger, error) {
	return core.NewLogger(cfg)
}
