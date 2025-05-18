package main

import (
	"context"
	"errors"
	"time"

	"github.com/gozephyr/cbreak"
	"github.com/gozephyr/examples/pkg/logger"
)

// CustomError represents a specific error type
type CustomError struct {
	Code    int
	Message string
}

func (e CustomError) Error() string {
	return e.Message
}

// IsTemporaryError determines if an error is temporary based on its code
func IsTemporaryError(err error) bool {
	var customErr CustomError
	if errors.As(err, &customErr) {
		// Consider 5xx errors as temporary
		return customErr.Code >= 500 && customErr.Code < 600
	}
	return false
}

func main() {
	log := logger.Get()
	log.SetPrefix("cbreak-advanced ")
	log.Section("Advanced Circuit Breaker Example")
	customFailureExample(log)
}

func customFailureExample(log *logger.Logger) {
	// Create a circuit breaker with custom failure detection
	config := cbreak.DefaultConfig("custom-failure-example")
	config.FailureThreshold = 3
	config.Timeout = 5 * time.Second
	config.CommandTimeout = 2 * time.Second
	config.HalfOpenMaxRequests = 1

	breaker, err := cbreak.NewBreaker[string](config)
	if err != nil {
		log.Error("Error creating circuit breaker: %v", err)
		return
	}
	defer breaker.Shutdown()

	// Simulate various error types
	errors := []error{
		CustomError{Code: 500, Message: "Internal Server Error"},
		CustomError{Code: 503, Message: "Service Unavailable"},
		CustomError{Code: 400, Message: "Bad Request"},
		CustomError{Code: 502, Message: "Bad Gateway"},
		CustomError{Code: 404, Message: "Not Found"},
	}

	log.Info("Simulating operations with different error types...")
	for i, err := range errors {
		result, execErr := breaker.Execute(context.Background(), func() (string, error) {
			return "", err
		})

		if execErr != nil {
			if IsTemporaryError(err) {
				log.Warn("Operation %d failed with temporary error: %v", i+1, err)
			} else {
				log.Error("Operation %d failed with permanent error: %v", i+1, err)
			}
		} else {
			log.Success("Operation %d succeeded with result: %s", i+1, result)
		}

		// Log the current state
		state := breaker.GetState()
		log.Info("Circuit breaker state: %s", state)
	}

	// Wait for the circuit breaker to reset
	log.Info("Waiting for circuit breaker to reset...")
	time.Sleep(6 * time.Second)

	// Try a successful operation
	log.Info("Trying a successful operation...")
	result, err := breaker.Execute(context.Background(), func() (string, error) {
		return "success", nil
	})

	if err != nil {
		log.Error("Operation failed: %v", err)
	} else {
		log.Success("Operation succeeded with result: %s", result)
	}

	// Log the final state
	state := breaker.GetState()
	log.Info("Final circuit breaker state: %s", state)
}
