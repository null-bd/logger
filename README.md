# Go Structured Logger

A lightweight, structured logging library for Go applications with support for multiple outputs, log levels, and field-based logging.

## Features

- Structured JSON logging
- Multiple output destinations (stdout, file)
- Log levels (Debug, Info, Warn, Error, Fatal)
- Context-aware request ID tracking
- Field-based logging
- Thread-safe operations
- Configurable through YAML/JSON
- Support for default fields
- Clean and simple API

## Installation

```bash
go get github.com/yourorg/logger
```

## Quick Start

```go
package main

import (
    "context"
    "github.com/yourorg/logger"
)

func main() {
    // Create a new logger
    log, err := logger.New(&logger.Config{
        ServiceName: "my-service",
        Environment: "development",
        LogLevel:    logger.InfoLevel,
        Format:      "json",
        OutputPaths: []string{"stdout"},
    })
    if err != nil {
        panic(err)
    }

    // Basic logging
    log.Info(context.Background(), "Application started", logger.Fields{
        "version": "1.0.0",
    })
}
```

## Configuration

The logger can be configured using the `Config` struct:

```go
type Config struct {
    ServiceName    string            `json:"service_name" yaml:"service_name"`
    Environment    string            `json:"environment" yaml:"environment"`
    LogLevel      Level             `json:"log_level" yaml:"log_level"`
    Format        string            `json:"format" yaml:"format"`
    DefaultFields map[string]string `json:"default_fields" yaml:"default_fields"`
    OutputPaths   []string          `json:"output_paths" yaml:"output_paths"`
}
```

### Configuration Options

- `ServiceName`: Name of your service (used in log entries)
- `Environment`: Environment name (e.g., "development", "production")
- `LogLevel`: Minimum log level to output ("debug", "info", "warn", "error", "fatal")
- `Format`: Log format ("json" or "text")
- `DefaultFields`: Fields to include in every log entry
- `OutputPaths`: Where to write logs (e.g., "stdout", "stderr", "/var/log/app.log")

## Usage Examples

### Basic Logging

```go
log.Debug(ctx, "Debug message", nil)
log.Info(ctx, "Info message", logger.Fields{"user_id": "123"})
log.Warn(ctx, "Warning message", logger.Fields{"latency": "100ms"})
log.Error(ctx, "Error occurred", logger.Fields{"error": err.Error()})
log.Fatal(ctx, "Fatal error", logger.Fields{"code": 500})
```

### With Default Fields

```go
config := &logger.Config{
    ServiceName: "user-service",
    Environment: "production",
    LogLevel:    logger.InfoLevel,
    Format:      "json",
    DefaultFields: map[string]string{
        "version": "1.0.0",
        "region":  "us-west-2",
    },
}

log, _ := logger.New(config)
```

### Multiple Outputs

```go
config := &logger.Config{
    ServiceName: "user-service",
    OutputPaths: []string{
        "stdout",
        "/var/log/app.log",
    },
}
```

### With Fields

```go
// Create a logger with additional fields
userLogger := log.WithFields(logger.Fields{
    "user_id": "123",
    "tenant":  "acme-corp",
})

// All logs will include the above fields
userLogger.Info(ctx, "User action completed", logger.Fields{
    "action": "profile_update",
})
```

## Output Examples

### JSON Format

```json
{
  "timestamp": "2024-12-15T10:30:45Z",
  "level": "info",
  "message": "User logged in",
  "service": "user-service",
  "environment": "production",
  "request_id": "5f6b7c8d-9e0f-1a2b-3c4d-5e6f7a8b9c0d",
  "fields": {
    "user_id": "123",
    "ip_address": "192.168.1.1",
    "version": "1.0.0"
  }
}
```

### Text Format

```
[2024-12-15T10:30:45Z] INFO: User logged in (RequestID: 5f6b7c8d-9e0f-1a2b-3c4d-5e6f7a8b9c0d)
```

## Best Practices

1. **Log Levels**: Use appropriate log levels:
   - `Debug`: Detailed information for debugging
   - `Info`: General operational entries
   - `Warn`: Warning messages for potentially harmful situations
   - `Error`: Error events that might still allow the application to continue running
   - `Fatal`: Severe errors that prevent the application from running

2. **Structured Fields**: Use structured fields instead of embedding data in messages:
   ```go
   // Good
   log.Info(ctx, "User created", logger.Fields{"user_id": "123"})
   
   // Avoid
   log.Info(ctx, "User 123 created", nil)
   ```

3. **Context**: Always pass a context to maintain request traceability:
   ```go
   ctx := context.Background()
   log.Info(ctx, "Processing started", fields)
   ```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.