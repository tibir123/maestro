package logger

import (
	"context"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

// Logger defines the interface for structured logging
type Logger interface {
	// Standard logging methods
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Panic(msg string, fields ...Field)

	// Context-aware logging
	WithContext(ctx context.Context) Logger
	WithFields(fields ...Field) Logger
	WithComponent(component string) Logger
	WithOperation(operation string) Logger

	// Utility methods
	IsDebugEnabled() bool
	GetLevel() string
}

// Field represents a key-value pair for structured logging
type Field struct {
	Key   string
	Value interface{}
}

// Fields is a slice of Field for convenience
type Fields []Field

// maestroLogger implements the Logger interface using logrus
type maestroLogger struct {
	entry     *logrus.Entry
	config    *Config
	mu        sync.RWMutex
	component string
}

var (
	// Global logger instance
	globalLogger Logger
	globalMutex  sync.RWMutex
)

// Initialize sets up the global logger with the provided configuration
func Initialize(config *Config) error {
	if config == nil {
		config = DefaultConfig()
	}

	if err := config.Validate(); err != nil {
		return err
	}

	logger, err := NewLogger(config)
	if err != nil {
		return err
	}

	globalMutex.Lock()
	globalLogger = logger
	globalMutex.Unlock()

	return nil
}

// InitializeDefault initializes the global logger with default configuration
func InitializeDefault() error {
	return Initialize(DefaultConfig())
}

// InitializeDevelopment initializes the global logger for development
func InitializeDevelopment() error {
	return Initialize(DevelopmentConfig())
}

// InitializeProduction initializes the global logger for production
func InitializeProduction() error {
	return Initialize(ProductionConfig())
}

// NewLogger creates a new logger instance with the given configuration
func NewLogger(config *Config) (Logger, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Create logrus instance
	log := logrus.New()

	// Set log level
	level, err := config.ParseLevel()
	if err != nil {
		return nil, err
	}
	log.SetLevel(level)

	// Set formatter
	if err := setFormatter(log, config); err != nil {
		return nil, err
	}

	// Set output
	if err := setOutput(log, config); err != nil {
		return nil, err
	}

	// Configure caller reporting
	log.SetReportCaller(config.EnableCaller)

	// Create base entry with component
	entry := log.WithField("component", config.Component)

	return &maestroLogger{
		entry:     entry,
		config:    config,
		component: config.Component,
	}, nil
}

// GetGlobal returns the global logger instance
func GetGlobal() Logger {
	globalMutex.RLock()
	defer globalMutex.RUnlock()

	if globalLogger == nil {
		// Initialize with default config if not initialized
		_ = InitializeDefault()
	}

	return globalLogger
}

// Component returns a logger instance configured for a specific component
func Component(component string) Logger {
	logger := GetGlobal()
	return logger.WithComponent(component)
}

// setFormatter configures the log formatter based on config
func setFormatter(log *logrus.Logger, config *Config) error {
	format := strings.ToLower(config.Format)

	switch format {
	case "json":
		formatter := &logrus.JSONFormatter{
			DisableTimestamp: !config.EnableTimestamp,
		}
		if config.EnableTimestamp && config.TimestampFormat != "" {
			formatter.TimestampFormat = config.TimestampFormat
		}
		log.SetFormatter(formatter)

	case "text":
		formatter := &logrus.TextFormatter{
			DisableTimestamp: !config.EnableTimestamp,
			FullTimestamp:    config.EnableTimestamp,
		}
		if config.EnableTimestamp && config.TimestampFormat != "" {
			formatter.TimestampFormat = config.TimestampFormat
		}
		log.SetFormatter(formatter)

	default:
		return NewConfigError("unsupported format: " + config.Format)
	}

	return nil
}

// setOutput configures the log output based on config
func setOutput(log *logrus.Logger, config *Config) error {
	output := strings.ToLower(config.Output)

	switch output {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "stderr":
		log.SetOutput(os.Stderr)
	default:
		// Assume it's a file path
		file, err := os.OpenFile(config.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		log.SetOutput(file)
	}

	return nil
}

// Implementation of Logger interface

func (l *maestroLogger) Debug(msg string, fields ...Field) {
	l.entry.WithFields(fieldsToLogrus(fields...)).Debug(msg)
}

func (l *maestroLogger) Info(msg string, fields ...Field) {
	l.entry.WithFields(fieldsToLogrus(fields...)).Info(msg)
}

func (l *maestroLogger) Warn(msg string, fields ...Field) {
	l.entry.WithFields(fieldsToLogrus(fields...)).Warn(msg)
}

func (l *maestroLogger) Error(msg string, fields ...Field) {
	l.entry.WithFields(fieldsToLogrus(fields...)).Error(msg)
}

func (l *maestroLogger) Fatal(msg string, fields ...Field) {
	l.entry.WithFields(fieldsToLogrus(fields...)).Fatal(msg)
}

func (l *maestroLogger) Panic(msg string, fields ...Field) {
	l.entry.WithFields(fieldsToLogrus(fields...)).Panic(msg)
}

func (l *maestroLogger) WithContext(ctx context.Context) Logger {
	// Extract common context values
	entry := l.entry

	if requestID := ctx.Value("request_id"); requestID != nil {
		entry = entry.WithField("request_id", requestID)
	}

	if sessionID := ctx.Value("session_id"); sessionID != nil {
		entry = entry.WithField("session_id", sessionID)
	}

	if userID := ctx.Value("user_id"); userID != nil {
		entry = entry.WithField("user_id", userID)
	}

	return &maestroLogger{
		entry:     entry,
		config:    l.config,
		component: l.component,
	}
}

func (l *maestroLogger) WithFields(fields ...Field) Logger {
	return &maestroLogger{
		entry:     l.entry.WithFields(fieldsToLogrus(fields...)),
		config:    l.config,
		component: l.component,
	}
}

func (l *maestroLogger) WithComponent(component string) Logger {
	return &maestroLogger{
		entry:     l.entry.WithField("component", component),
		config:    l.config,
		component: component,
	}
}

func (l *maestroLogger) WithOperation(operation string) Logger {
	return &maestroLogger{
		entry:     l.entry.WithField("operation", operation),
		config:    l.config,
		component: l.component,
	}
}

func (l *maestroLogger) IsDebugEnabled() bool {
	return l.entry.Logger.IsLevelEnabled(logrus.DebugLevel)
}

func (l *maestroLogger) GetLevel() string {
	return l.entry.Logger.GetLevel().String()
}

// Helper functions

// fieldsToLogrus converts our Field slice to logrus.Fields
func fieldsToLogrus(fields ...Field) logrus.Fields {
	if len(fields) == 0 {
		return nil
	}

	logrusFields := make(logrus.Fields, len(fields))
	for _, field := range fields {
		logrusFields[field.Key] = field.Value
	}
	return logrusFields
}

// Convenience functions for creating fields
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

func Float64(key string, value float64) Field {
	return Field{Key: key, Value: value}
}

func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

func Error(err error) Field {
	return Field{Key: "error", Value: err.Error()}
}

func Duration(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// Global convenience functions that use the global logger

// Debug logs a debug message using the global logger
func Debug(msg string, fields ...Field) {
	GetGlobal().Debug(msg, fields...)
}

// Info logs an info message using the global logger
func Info(msg string, fields ...Field) {
	GetGlobal().Info(msg, fields...)
}

// Warn logs a warning message using the global logger
func Warn(msg string, fields ...Field) {
	GetGlobal().Warn(msg, fields...)
}

// ErrorMsg logs an error message using the global logger
func ErrorMsg(msg string, fields ...Field) {
	GetGlobal().Error(msg, fields...)
}

// Fatal logs a fatal message using the global logger and exits
func Fatal(msg string, fields ...Field) {
	GetGlobal().Fatal(msg, fields...)
}

// Panic logs a panic message using the global logger and panics
func Panic(msg string, fields ...Field) {
	GetGlobal().Panic(msg, fields...)
}

// WithContext returns a logger with context using the global logger
func WithContext(ctx context.Context) Logger {
	return GetGlobal().WithContext(ctx)
}

// WithFields returns a logger with fields using the global logger
func WithFields(fields ...Field) Logger {
	return GetGlobal().WithFields(fields...)
}

// WithOperation returns a logger with operation using the global logger
func WithOperation(operation string) Logger {
	return GetGlobal().WithOperation(operation)
}

// SetGlobalOutput allows changing the output destination at runtime
func SetGlobalOutput(output io.Writer) {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	if globalLogger != nil {
		if ml, ok := globalLogger.(*maestroLogger); ok {
			ml.entry.Logger.SetOutput(output)
		}
	}
}
