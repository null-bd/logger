package types

// Level represents log severity
type Level string

const (
	DebugLevel Level = "debug"
	InfoLevel  Level = "info"
	WarnLevel  Level = "warn"
	ErrorLevel Level = "error"
	FatalLevel Level = "fatal"
)

// Fields represents structured log fields
type Fields map[string]interface{}

// Logger defines the interface for logging operations
type Logger interface {
	Debug(msg string, fields Fields)
	Info(msg string, fields Fields)
	Warn(msg string, fields Fields)
	Error(msg string, fields Fields)
	Fatal(msg string, fields Fields)
	WithFields(fields Fields) Logger
}

// Config holds the logger configuration
type Config struct {
	ServiceName   string            `json:"service_name" yaml:"service_name"`
	Environment   string            `json:"environment" yaml:"environment"`
	LogLevel      Level             `json:"log_level" yaml:"log_level"`
	Format        string            `json:"format" yaml:"format"`
	DefaultFields map[string]string `json:"default_fields" yaml:"default_fields"`
	OutputPaths   []string          `json:"output_paths" yaml:"output_paths"`
}
