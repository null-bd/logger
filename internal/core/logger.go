package core

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/null-bd/logger/types"
)

type logEntry struct {
	Timestamp   time.Time    `json:"timestamp"`
	Level       types.Level  `json:"level"`
	Message     string       `json:"message"`
	Service     string       `json:"service"`
	Environment string       `json:"environment"`
	RequestID   string       `json:"request_id,omitempty"`
	TraceID     string       `json:"trace_id,omitempty"`
	Fields      types.Fields `json:"fields,omitempty"`
}

type Logger struct {
	mu            sync.RWMutex
	config        *types.Config
	writers       []io.Writer
	defaultFields types.Fields
}

func NewLogger(cfg *types.Config) (types.Logger, error) {
	if cfg == nil {
		cfg = defaultConfig()
	}

	l := &Logger{
		config:        cfg,
		defaultFields: make(types.Fields),
		writers:       make([]io.Writer, 0),
	}

	if err := l.initializeWriters(); err != nil {
		return nil, err
	}

	for k, v := range cfg.DefaultFields {
		l.defaultFields[k] = v
	}

	return l, nil
}

func (l *Logger) initializeWriters() error {
	for _, path := range l.config.OutputPaths {
		writer, err := l.createWriter(path)
		if err != nil {
			return fmt.Errorf("failed to create writer for %s: %v", path, err)
		}
		l.writers = append(l.writers, writer)
	}
	return nil
}

func (l *Logger) createWriter(path string) (io.Writer, error) {
	switch path {
	case "stdout":
		return os.Stdout, nil
	case "stderr":
		return os.Stderr, nil
	default:
		return os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}
}

func (l *Logger) log(level types.Level, msg string, fields types.Fields) {
	if !l.isLevelEnabled(level) {
		return
	}

	entry := &logEntry{
		Timestamp:   time.Now().UTC(),
		Level:       level,
		Message:     msg,
		Service:     l.config.ServiceName,
		Environment: l.config.Environment,
		Fields:      make(types.Fields),
	}

	for k, v := range GetTraceFields() {
		switch k {
		case "request_id":
			entry.RequestID = v.(string)
		case "trace_id":
			entry.TraceID = v.(string)
		default:
			entry.Fields[k] = v
		}
	}

	l.mergeFields(entry, fields)
	l.writeLog(entry)
}

func (l *Logger) isLevelEnabled(level types.Level) bool {
	levels := map[types.Level]int{
		types.DebugLevel: 0,
		types.InfoLevel:  1,
		types.WarnLevel:  2,
		types.ErrorLevel: 3,
		types.FatalLevel: 4,
	}
	return levels[level] >= levels[l.config.LogLevel]
}

func (l *Logger) mergeFields(entry *logEntry, fields types.Fields) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	for k, v := range l.defaultFields {
		entry.Fields[k] = v
	}

	for k, v := range fields {
		entry.Fields[k] = v
	}
}

func (l *Logger) writeLog(entry *logEntry) {
	var output []byte
	var err error

	if l.config.Format == "json" {
		output, err = json.Marshal(entry)
	} else {
		output = []byte(fmt.Sprintf("[%s] %s: %s (RequestID: %s)\n",
			entry.Timestamp.Format(time.RFC3339),
			entry.Level,
			entry.Message,
			entry.RequestID))
	}

	if err != nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	for _, w := range l.writers {
		w.Write(output)
		if l.config.Format == "json" {
			w.Write([]byte("\n"))
		}
	}
}

// Logger interface implementation
func (l *Logger) Debug(msg string, fields types.Fields) {
	l.log(types.DebugLevel, msg, fields)
}

func (l *Logger) Info(msg string, fields types.Fields) {
	l.log(types.InfoLevel, msg, fields)
}

func (l *Logger) Warn(msg string, fields types.Fields) {
	l.log(types.WarnLevel, msg, fields)
}

func (l *Logger) Error(msg string, fields types.Fields) {
	l.log(types.ErrorLevel, msg, fields)
}

func (l *Logger) Fatal(msg string, fields types.Fields) {
	l.log(types.FatalLevel, msg, fields)
	os.Exit(1)
}

func (l *Logger) WithFields(fields types.Fields) types.Logger {
	newLogger := &Logger{
		config:        l.config,
		writers:       l.writers,
		defaultFields: make(types.Fields),
	}

	for k, v := range l.defaultFields {
		newLogger.defaultFields[k] = v
	}

	for k, v := range fields {
		newLogger.defaultFields[k] = v
	}

	return newLogger
}
