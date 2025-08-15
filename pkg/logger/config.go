package logger

import (
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Config defines the configuration for the logging system
type Config struct {
	// Level defines the minimum log level (debug, info, warn, error)
	Level string `mapstructure:"level" json:"level"`

	// Format specifies the log format (json, text)
	Format string `mapstructure:"format" json:"format"`

	// Output specifies where to write logs (stdout, stderr, file path)
	Output string `mapstructure:"output" json:"output"`

	// Component is the component name that will be included in all logs
	Component string `mapstructure:"component" json:"component"`

	// EnableCaller adds caller information to logs
	EnableCaller bool `mapstructure:"enable_caller" json:"enable_caller"`

	// EnableTimestamp adds timestamp to logs
	EnableTimestamp bool `mapstructure:"enable_timestamp" json:"enable_timestamp"`

	// TimestampFormat defines the timestamp format
	TimestampFormat string `mapstructure:"timestamp_format" json:"timestamp_format"`

	// FileRotation configures log file rotation (only used when Output is a file)
	FileRotation *FileRotationConfig `mapstructure:"file_rotation" json:"file_rotation,omitempty"`
}

// FileRotationConfig defines configuration for log file rotation
type FileRotationConfig struct {
	// MaxSize is the maximum size in MB before rotation
	MaxSize int `mapstructure:"max_size" json:"max_size"`

	// MaxAge is the maximum age in days to keep old log files
	MaxAge int `mapstructure:"max_age" json:"max_age"`

	// MaxBackups is the maximum number of old log files to keep
	MaxBackups int `mapstructure:"max_backups" json:"max_backups"`

	// Compress determines if old log files should be compressed
	Compress bool `mapstructure:"compress" json:"compress"`
}

// DefaultConfig returns a default logging configuration
func DefaultConfig() *Config {
	return &Config{
		Level:           "info",
		Format:          "json",
		Output:          "stdout",
		Component:       "maestro",
		EnableCaller:    false,
		EnableTimestamp: true,
		TimestampFormat: time.RFC3339,
		FileRotation: &FileRotationConfig{
			MaxSize:    100, // 100MB
			MaxAge:     30,  // 30 days
			MaxBackups: 5,   // 5 backup files
			Compress:   true,
		},
	}
}

// DevelopmentConfig returns a configuration suitable for development
func DevelopmentConfig() *Config {
	config := DefaultConfig()
	config.Level = "debug"
	config.Format = "text"
	config.EnableCaller = true
	return config
}

// ProductionConfig returns a configuration suitable for production
func ProductionConfig() *Config {
	config := DefaultConfig()
	config.Level = "info"
	config.Format = "json"
	config.EnableCaller = false
	return config
}

// ComponentConfig returns a new config with the specified component name
func (c *Config) ComponentConfig(component string) *Config {
	newConfig := *c // Copy the config
	newConfig.Component = component
	return &newConfig
}

// ParseLevel converts a string level to logrus.Level
func (c *Config) ParseLevel() (logrus.Level, error) {
	return logrus.ParseLevel(strings.ToLower(c.Level))
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Validate log level
	if _, err := c.ParseLevel(); err != nil {
		return err
	}

	// Validate format
	format := strings.ToLower(c.Format)
	if format != "json" && format != "text" {
		return NewConfigError("format must be 'json' or 'text'")
	}

	// Validate output
	if c.Output == "" {
		return NewConfigError("output cannot be empty")
	}

	// Validate timestamp format
	if c.EnableTimestamp && c.TimestampFormat == "" {
		return NewConfigError("timestamp_format cannot be empty when enable_timestamp is true")
	}

	return nil
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() *Config {
	config := DefaultConfig()

	if level := os.Getenv("MAESTRO_LOG_LEVEL"); level != "" {
		config.Level = level
	}

	if format := os.Getenv("MAESTRO_LOG_FORMAT"); format != "" {
		config.Format = format
	}

	if output := os.Getenv("MAESTRO_LOG_OUTPUT"); output != "" {
		config.Output = output
	}

	if component := os.Getenv("MAESTRO_LOG_COMPONENT"); component != "" {
		config.Component = component
	}

	if caller := os.Getenv("MAESTRO_LOG_CALLER"); caller == "true" {
		config.EnableCaller = true
	}

	return config
}

// ConfigError represents a configuration validation error
type ConfigError struct {
	Message string
}

func (e *ConfigError) Error() string {
	return "logger config error: " + e.Message
}

// NewConfigError creates a new configuration error
func NewConfigError(message string) *ConfigError {
	return &ConfigError{Message: message}
}
