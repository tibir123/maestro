package logger

import (
	"context"
	"fmt"
	"time"

	"github.com/madstone-tech/maestro/domain/music"
)

// This file demonstrates how to integrate logging into infrastructure components

// ExampleAppleScriptExecutorWithLogging shows how to add comprehensive logging to the AppleScript executor
type ExampleAppleScriptExecutor struct {
	logger  Logger
	timeout time.Duration
}

func NewExampleAppleScriptExecutor(timeout time.Duration) *ExampleAppleScriptExecutor {
	return &ExampleAppleScriptExecutor{
		logger:  Component("applescript"),
		timeout: timeout,
	}
}

// ExecuteScript demonstrates logging in AppleScript execution
func (e *ExampleAppleScriptExecutor) ExecuteScript(ctx context.Context, scriptName string, args ...string) (string, error) {
	logger := e.logger.WithOperation("execute_script").WithContext(ctx)

	start := time.Now()
	logger.Debug("Starting AppleScript execution",
		String("script", scriptName),
		String("args", fmt.Sprintf("%v", args)),
		Duration("timeout", e.timeout),
	)

	// Simulate script execution with timeout
	scriptCtx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	// Monitor for context cancellation
	done := make(chan struct {
		result string
		err    error
	}, 1)

	go func() {
		// Simulate script execution
		time.Sleep(100 * time.Millisecond) // Simulate work

		// Simulate occasional failures
		if time.Now().UnixNano()%10 == 0 {
			done <- struct {
				result string
				err    error
			}{"", fmt.Errorf("osascript execution failed")}
			return
		}

		done <- struct {
			result string
			err    error
		}{"script output", nil}
	}()

	select {
	case result := <-done:
		duration := time.Since(start)

		if result.err != nil {
			logger.Error("AppleScript execution failed",
				Error(result.err),
				String("script", scriptName),
				Duration("execution_time", duration),
			)
			return "", result.err
		}

		logger.Info("AppleScript execution completed",
			String("script", scriptName),
			Duration("execution_time", duration),
			Int("output_length", len(result.result)),
		)

		// Log slow executions
		if duration > 500*time.Millisecond {
			logger.Warn("Slow AppleScript execution detected",
				String("script", scriptName),
				Duration("execution_time", duration),
				Duration("threshold", 500*time.Millisecond),
			)
		}

		return result.result, nil

	case <-scriptCtx.Done():
		duration := time.Since(start)
		logger.Error("AppleScript execution timed out",
			String("script", scriptName),
			Duration("timeout", e.timeout),
			Duration("elapsed_time", duration),
		)
		return "", fmt.Errorf("script execution timed out after %v", e.timeout)
	}
}

// ExamplePlayerRepositoryWithLogging demonstrates logging in the player repository
type ExamplePlayerRepository struct {
	executor *ExampleAppleScriptExecutor
	logger   Logger
}

func NewExamplePlayerRepository(executor *ExampleAppleScriptExecutor) *ExamplePlayerRepository {
	return &ExamplePlayerRepository{
		executor: executor,
		logger:   Component("player_repository"),
	}
}

func (r *ExamplePlayerRepository) Play(ctx context.Context) error {
	logger := r.logger.WithOperation("play").WithContext(ctx)

	start := time.Now()
	logger.Debug("Executing play command")

	output, err := r.executor.ExecuteScript(ctx, "play.scpt")
	if err != nil {
		logger.Error("Play command failed",
			Error(err),
			Duration("operation_time", time.Since(start)),
		)
		return err
	}

	logger.Info("Play command succeeded",
		Duration("operation_time", time.Since(start)),
		String("script_output", output),
	)

	return nil
}

func (r *ExamplePlayerRepository) GetCurrentState(ctx context.Context) (*music.Player, error) {
	logger := r.logger.WithOperation("get_current_state").WithContext(ctx)

	start := time.Now()
	logger.Debug("Getting current player state")

	_, err := r.executor.ExecuteScript(ctx, "get_player_state.scpt")
	if err != nil {
		logger.Error("Failed to get player state",
			Error(err),
			Duration("operation_time", time.Since(start)),
		)
		return nil, err
	}

	// Simulate parsing the output
	player := &music.Player{
		State:    music.PlayerStatePlaying,
		Volume:   music.NewVolume(75),
		Position: music.NewDuration(120),
		Shuffle:  false,
		Repeat:   music.RepeatModeOff,
	}

	logger.Info("Retrieved player state",
		String("state", player.State.String()),
		Int("volume", player.Volume.Level()),
		String("position", player.Position.String()),
		Duration("operation_time", time.Since(start)),
	)

	return player, nil
}

// ExampleSessionManagerWithLogging demonstrates logging in session management
type ExampleSessionManager struct {
	sessions map[string]*Session
	logger   Logger
}

type Session struct {
	ID        string
	ClientID  string
	Type      string
	CreatedAt time.Time
	LastSeen  time.Time
}

func NewExampleSessionManager() *ExampleSessionManager {
	return &ExampleSessionManager{
		sessions: make(map[string]*Session),
		logger:   Component("session_manager"),
	}
}

func (sm *ExampleSessionManager) CreateSession(ctx context.Context, clientID, sessionType string) (*Session, error) {
	logger := sm.logger.WithOperation("create_session").WithContext(ctx)

	sessionID := fmt.Sprintf("sess_%d", time.Now().UnixNano())
	session := &Session{
		ID:        sessionID,
		ClientID:  clientID,
		Type:      sessionType,
		CreatedAt: time.Now(),
		LastSeen:  time.Now(),
	}

	// Check for existing sessions
	existingCount := len(sm.sessions)
	if existingCount >= 10 {
		logger.Warn("Maximum concurrent sessions reached",
			Int("current_sessions", existingCount),
			Int("max_sessions", 10),
		)
		return nil, fmt.Errorf("maximum sessions exceeded")
	}

	sm.sessions[sessionID] = session

	logger.Info("Session created",
		String("session_id", sessionID),
		String("client_id", clientID),
		String("session_type", sessionType),
		Int("total_sessions", len(sm.sessions)),
	)

	return session, nil
}

func (sm *ExampleSessionManager) CleanupExpiredSessions(ctx context.Context) {
	logger := sm.logger.WithOperation("cleanup_expired_sessions").WithContext(ctx)

	start := time.Now()
	timeout := 5 * time.Minute
	expiredCount := 0

	logger.Debug("Starting session cleanup",
		Duration("session_timeout", timeout),
		Int("current_sessions", len(sm.sessions)),
	)

	for sessionID, session := range sm.sessions {
		if time.Since(session.LastSeen) > timeout {
			delete(sm.sessions, sessionID)
			expiredCount++

			logger.Debug("Session expired and removed",
				String("session_id", sessionID),
				String("client_id", session.ClientID),
				Duration("idle_time", time.Since(session.LastSeen)),
			)
		}
	}

	logger.Info("Session cleanup completed",
		Int("expired_sessions", expiredCount),
		Int("remaining_sessions", len(sm.sessions)),
		Duration("cleanup_time", time.Since(start)),
	)
}

// ExampleGRPCServerWithLogging demonstrates logging in gRPC server
type ExampleGRPCServer struct {
	logger Logger
	port   int
}

func NewExampleGRPCServer(port int) *ExampleGRPCServer {
	return &ExampleGRPCServer{
		logger: Component("grpc_server"),
		port:   port,
	}
}

func (s *ExampleGRPCServer) Start(ctx context.Context) error {
	logger := s.logger.WithOperation("start_server").WithContext(ctx)

	logger.Info("Starting gRPC server",
		Int("port", s.port),
		Bool("tls_enabled", true),
	)

	// Simulate server startup
	time.Sleep(100 * time.Millisecond)

	logger.Info("gRPC server started successfully",
		String("address", fmt.Sprintf(":%d", s.port)),
		String("status", "ready"),
	)

	return nil
}

func (s *ExampleGRPCServer) HandlePlayRequest(ctx context.Context, sessionID string) error {
	logger := s.logger.WithOperation("handle_play_request").WithContext(ctx)

	start := time.Now()
	logger.Debug("Processing play request",
		String("session_id", sessionID),
	)

	// Simulate request processing
	time.Sleep(50 * time.Millisecond)

	logger.Info("Play request processed",
		String("session_id", sessionID),
		Duration("processing_time", time.Since(start)),
		Bool("success", true),
	)

	return nil
}

// ExampleCacheWithLogging demonstrates logging in caching layer
type ExampleCache struct {
	data   map[string]interface{}
	logger Logger
}

func NewExampleCache() *ExampleCache {
	return &ExampleCache{
		data:   make(map[string]interface{}),
		logger: Component("cache"),
	}
}

func (c *ExampleCache) Get(ctx context.Context, key string) (interface{}, bool) {
	logger := c.logger.WithOperation("cache_get").WithContext(ctx)

	start := time.Now()
	value, exists := c.data[key]

	logger.Debug("Cache lookup",
		String("key", key),
		Bool("hit", exists),
		Duration("lookup_time", time.Since(start)),
	)

	if exists {
		logger.Debug("Cache hit",
			String("key", key),
			String("value_type", fmt.Sprintf("%T", value)),
		)
	}

	return value, exists
}

func (c *ExampleCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) {
	logger := c.logger.WithOperation("cache_set").WithContext(ctx)

	start := time.Now()
	c.data[key] = value

	logger.Debug("Cache entry stored",
		String("key", key),
		String("value_type", fmt.Sprintf("%T", value)),
		Duration("ttl", ttl),
		Duration("store_time", time.Since(start)),
		Int("cache_size", len(c.data)),
	)
}

// ExampleHealthCheckerWithLogging demonstrates logging in health checks
type ExampleHealthChecker struct {
	logger     Logger
	playerRepo *ExamplePlayerRepository
}

func NewExampleHealthChecker(playerRepo *ExamplePlayerRepository) *ExampleHealthChecker {
	return &ExampleHealthChecker{
		logger:     Component("health_checker"),
		playerRepo: playerRepo,
	}
}

func (h *ExampleHealthChecker) CheckHealth(ctx context.Context) map[string]bool {
	logger := h.logger.WithOperation("health_check").WithContext(ctx)

	start := time.Now()
	logger.Debug("Starting health check")

	results := make(map[string]bool)

	// Check Music.app connectivity
	_, err := h.playerRepo.GetCurrentState(ctx)
	musicAppHealthy := err == nil
	results["music_app"] = musicAppHealthy

	if musicAppHealthy {
		logger.Debug("Music.app health check passed")
	} else {
		logger.Warn("Music.app health check failed", Error(err))
	}

	// Check system resources (simulated)
	memoryOK := true // Simulate memory check
	results["memory"] = memoryOK

	// Check disk space (simulated)
	diskSpaceOK := true // Simulate disk check
	results["disk_space"] = diskSpaceOK

	overall := musicAppHealthy && memoryOK && diskSpaceOK
	results["overall"] = overall

	logger.Info("Health check completed",
		Bool("music_app_healthy", musicAppHealthy),
		Bool("memory_ok", memoryOK),
		Bool("disk_space_ok", diskSpaceOK),
		Bool("overall_healthy", overall),
		Duration("check_time", time.Since(start)),
	)

	return results
}
