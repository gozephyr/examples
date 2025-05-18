package main

import (
	"context"
	"errors"
	"time"

	"github.com/gozephyr/cbreak"
	"github.com/gozephyr/examples/pkg/logger"
)

func main() {
	log := logger.Get()
	log.SetPrefix("cbreak-simple ")
	log.Section("Simple Circuit Breaker Example")
	simpleExample(log)
}

func simpleExample(log *logger.Logger) {
	// Create a circuit breaker with basic configuration
	config := cbreak.DefaultConfig("simple-example")
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

	// Simulate a failing operation
	log.Info("Simulating failing operations...")
	for i := 0; i < 5; i++ {
		result, err := breaker.Execute(context.Background(), func() (string, error) {
			// Simulate a failing operation
			return "", errors.New("operation failed")
		})

		if err != nil {
			log.Error("Operation %d failed: %v", i+1, err)
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
		// Simulate a successful operation
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
