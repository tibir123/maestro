package applescript

import (
	"context"
	"time"

	"github.com/madstone-tech/maestro/domain/music"
)

// DefaultConfig returns a default configuration suitable for most use cases.
func DefaultConfig() *ExecutorConfig {
	return DefaultExecutorConfig()
}

// NewDefaultPlayerRepository creates a PlayerRepository with default settings.
// This is a convenience function for the most common use case.
func NewDefaultPlayerRepository() music.PlayerRepository {
	executor := NewExecutor(nil)
	return NewPlayerRepository(executor)
}

// NewPlayerRepositoryWithConfig creates a PlayerRepository with custom configuration.
func NewPlayerRepositoryWithConfig(config *ExecutorConfig) music.PlayerRepository {
	executor := NewExecutor(config)
	return NewPlayerRepository(executor)
}

// QuickHealthCheck performs a fast health check to ensure the infrastructure is working.
// This is useful for application startup validation.
func QuickHealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	executor := NewExecutor(nil)
	if err := executor.IsExecutable(); err != nil {
		return err
	}

	return executor.HealthCheck(ctx)
}

// CheckMusicAppAvailability specifically checks if Music.app is accessible.
func CheckMusicAppAvailability() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	playerRepo := NewDefaultPlayerRepository()

	// Type assertion to access health check method
	if repo, ok := playerRepo.(*PlayerRepository); ok {
		return repo.HealthCheck(ctx)
	}

	// Fallback: try to get player state
	_, err := playerRepo.GetCurrentState(ctx)
	return err
}
