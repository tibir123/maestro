package logger

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.Level != "info" {
		t.Errorf("Expected default level 'info', got '%s'", config.Level)
	}

	if config.Format != "json" {
		t.Errorf("Expected default format 'json', got '%s'", config.Format)
	}

	if config.Component != "maestro" {
		t.Errorf("Expected default component 'maestro', got '%s'", config.Component)
	}

	if err := config.Validate(); err != nil {
		t.Errorf("Default config should be valid, got error: %v", err)
	}
}

func TestDevelopmentConfig(t *testing.T) {
	config := DevelopmentConfig()

	if config.Level != "debug" {
		t.Errorf("Expected development level 'debug', got '%s'", config.Level)
	}

	if config.Format != "text" {
		t.Errorf("Expected development format 'text', got '%s'", config.Format)
	}

	if !config.EnableCaller {
		t.Error("Expected development config to enable caller")
	}
}

func TestProductionConfig(t *testing.T) {
	config := ProductionConfig()

	if config.Level != "info" {
		t.Errorf("Expected production level 'info', got '%s'", config.Level)
	}

	if config.Format != "json" {
		t.Errorf("Expected production format 'json', got '%s'", config.Format)
	}

	if config.EnableCaller {
		t.Error("Expected production config to disable caller")
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		expectErr bool
	}{
		{
			name:      "valid config",
			config:    DefaultConfig(),
			expectErr: false,
		},
		{
			name: "invalid level",
			config: &Config{
				Level:           "invalid",
				Format:          "json",
				Output:          "stdout",
				Component:       "test",
				EnableTimestamp: true,
				TimestampFormat: time.RFC3339,
			},
			expectErr: true,
		},
		{
			name: "invalid format",
			config: &Config{
				Level:           "info",
				Format:          "invalid",
				Output:          "stdout",
				Component:       "test",
				EnableTimestamp: true,
				TimestampFormat: time.RFC3339,
			},
			expectErr: true,
		},
		{
			name: "empty output",
			config: &Config{
				Level:           "info",
				Format:          "json",
				Output:          "",
				Component:       "test",
				EnableTimestamp: true,
				TimestampFormat: time.RFC3339,
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectErr && err == nil {
				t.Error("Expected validation error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no validation error but got: %v", err)
			}
		})
	}
}

func TestLoadFromEnv(t *testing.T) {
	// Save original env vars
	originalVars := map[string]string{
		"MAESTRO_LOG_LEVEL":     os.Getenv("MAESTRO_LOG_LEVEL"),
		"MAESTRO_LOG_FORMAT":    os.Getenv("MAESTRO_LOG_FORMAT"),
		"MAESTRO_LOG_OUTPUT":    os.Getenv("MAESTRO_LOG_OUTPUT"),
		"MAESTRO_LOG_COMPONENT": os.Getenv("MAESTRO_LOG_COMPONENT"),
		"MAESTRO_LOG_CALLER":    os.Getenv("MAESTRO_LOG_CALLER"),
	}

	// Clean up after test
	defer func() {
		for key, value := range originalVars {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	// Set test environment variables
	os.Setenv("MAESTRO_LOG_LEVEL", "debug")
	os.Setenv("MAESTRO_LOG_FORMAT", "text")
	os.Setenv("MAESTRO_LOG_OUTPUT", "stderr")
	os.Setenv("MAESTRO_LOG_COMPONENT", "test")
	os.Setenv("MAESTRO_LOG_CALLER", "true")

	config := LoadFromEnv()

	if config.Level != "debug" {
		t.Errorf("Expected level 'debug', got '%s'", config.Level)
	}
	if config.Format != "text" {
		t.Errorf("Expected format 'text', got '%s'", config.Format)
	}
	if config.Output != "stderr" {
		t.Errorf("Expected output 'stderr', got '%s'", config.Output)
	}
	if config.Component != "test" {
		t.Errorf("Expected component 'test', got '%s'", config.Component)
	}
	if !config.EnableCaller {
		t.Error("Expected caller to be enabled")
	}
}

func TestNewLogger(t *testing.T) {
	config := &Config{
		Level:           "debug",
		Format:          "json",
		Output:          "stdout",
		Component:       "test",
		EnableCaller:    false,
		EnableTimestamp: true,
		TimestampFormat: time.RFC3339,
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	if logger == nil {
		t.Fatal("Expected logger to be created")
	}

	if !logger.IsDebugEnabled() {
		t.Error("Expected debug to be enabled")
	}

	if logger.GetLevel() != "debug" {
		t.Errorf("Expected level 'debug', got '%s'", logger.GetLevel())
	}
}

func TestLoggerMethods(t *testing.T) {
	config := DevelopmentConfig()
	config.Format = "json"

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Note: This is a simplified test. In real scenarios, you might need to
	// use more sophisticated output capture methods.

	// Test basic logging methods
	logger.Debug("debug message", String("key", "value"))
	logger.Info("info message", String("key", "value"))
	logger.Warn("warn message", String("key", "value"))
	logger.Error("error message", String("key", "value"))

	// Test logger with fields
	fieldLogger := logger.WithFields(
		String("component", "test"),
		Int("request_id", 123),
	)
	fieldLogger.Info("message with fields")

	// Test logger with component
	componentLogger := logger.WithComponent("test_component")
	componentLogger.Info("component message")

	// Test logger with operation
	operationLogger := logger.WithOperation("test_operation")
	operationLogger.Info("operation message")
}

func TestLoggerWithContext(t *testing.T) {
	config := DevelopmentConfig()
	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Create context with values
	ctx := context.Background()
	ctx = context.WithValue(ctx, "request_id", "req-123")
	ctx = context.WithValue(ctx, "session_id", "sess-456")
	ctx = context.WithValue(ctx, "user_id", "user-789")

	// Create context-aware logger
	contextLogger := logger.WithContext(ctx)
	contextLogger.Info("message with context")

	// The context values should be included in the log output
	// In a real test, you would capture and verify the output
}

func TestFieldHelpers(t *testing.T) {
	// Test basic field types
	stringField := String("key", "value")
	if stringField.Key != "key" || stringField.Value != "value" {
		t.Errorf("String field: expected key='key' value='value', got key='%s' value='%v'", stringField.Key, stringField.Value)
	}

	intField := Int("key", 42)
	if intField.Key != "key" || intField.Value != 42 {
		t.Errorf("Int field: expected key='key' value=42, got key='%s' value=%v", intField.Key, intField.Value)
	}

	int64Field := Int64("key", int64(42))
	if int64Field.Key != "key" || int64Field.Value != int64(42) {
		t.Errorf("Int64 field: expected key='key' value=42, got key='%s' value=%v", int64Field.Key, int64Field.Value)
	}

	floatField := Float64("key", 3.14)
	if floatField.Key != "key" || floatField.Value != 3.14 {
		t.Errorf("Float64 field: expected key='key' value=3.14, got key='%s' value=%v", floatField.Key, floatField.Value)
	}

	boolField := Bool("key", true)
	if boolField.Key != "key" || boolField.Value != true {
		t.Errorf("Bool field: expected key='key' value=true, got key='%s' value=%v", boolField.Key, boolField.Value)
	}

	durationField := Duration("key", 5*time.Second)
	if durationField.Key != "key" || durationField.Value != 5*time.Second {
		t.Errorf("Duration field: expected key='key' value=5s, got key='%s' value=%v", durationField.Key, durationField.Value)
	}

	// Test Any field (just verify it's set, don't compare complex types)
	anyField := Any("key", map[string]string{"test": "value"})
	if anyField.Key != "key" || anyField.Value == nil {
		t.Errorf("Any field: expected key='key' with non-nil value, got key='%s' value=%v", anyField.Key, anyField.Value)
	}
}

func TestErrorField(t *testing.T) {
	testErr := NewConfigError("test error")
	field := Error(testErr)

	if field.Key != "error" {
		t.Errorf("Expected error field key 'error', got '%s'", field.Key)
	}

	if field.Value != testErr.Error() {
		t.Errorf("Expected error field value '%s', got '%s'", testErr.Error(), field.Value)
	}
}

func TestGlobalLogger(t *testing.T) {
	// Test that global logger can be initialized
	err := InitializeDefault()
	if err != nil {
		t.Fatalf("Failed to initialize global logger: %v", err)
	}

	// Test that we can get the global logger
	logger := GetGlobal()
	if logger == nil {
		t.Fatal("Expected global logger to be available")
	}

	// Test global convenience functions
	Debug("test debug message")
	Info("test info message")
	Warn("test warn message")
	ErrorMsg("test error message")

	// Test component logger
	componentLogger := Component("test_component")
	if componentLogger == nil {
		t.Fatal("Expected component logger to be created")
	}
}

func TestComponentConfig(t *testing.T) {
	config := DefaultConfig()
	originalComponent := config.Component

	newConfig := config.ComponentConfig("new_component")

	// Original config should be unchanged
	if config.Component != originalComponent {
		t.Error("Original config was modified")
	}

	// New config should have different component
	if newConfig.Component != "new_component" {
		t.Errorf("Expected new component 'new_component', got '%s'", newConfig.Component)
	}
}

func TestJSONOutput(t *testing.T) {
	config := &Config{
		Level:           "info",
		Format:          "json",
		Output:          "stdout", // We'll override this for testing
		Component:       "test",
		EnableCaller:    false,
		EnableTimestamp: true,
		TimestampFormat: time.RFC3339,
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// For this test, we'll verify that the logger was created successfully
	// In a more comprehensive test, you would capture the actual output
	// and verify the JSON structure

	logger.Info("test message", String("field", "value"))

	// The logger should be functional without errors
	if !logger.IsDebugEnabled() {
		// Debug should be disabled for info level
	}
}

// Benchmark tests
func BenchmarkLogger(b *testing.B) {
	config := ProductionConfig()
	config.Output = "/dev/null" // Discard output for benchmarking

	logger, err := NewLogger(config)
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("benchmark message",
				String("key1", "value1"),
				Int("key2", 42),
				Bool("key3", true),
			)
		}
	})
}

func BenchmarkLoggerWithFields(b *testing.B) {
	config := ProductionConfig()
	config.Output = "/dev/null"

	logger, err := NewLogger(config)
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}

	fieldLogger := logger.WithFields(
		String("component", "benchmark"),
		String("version", "1.0.0"),
	)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			fieldLogger.Info("benchmark message",
				String("operation", "test"),
				Int("iteration", 1),
			)
		}
	})
}
