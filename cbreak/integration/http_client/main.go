package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gozephyr/cbreak"
	"github.com/gozephyr/examples/pkg/logger"
)

func main() {
	log := logger.Get()
	log.SetPrefix("cbreak-http ")
	log.Section("HTTP Client Integration Example")
	httpClientExample(log)
}

func httpClientExample(log *logger.Logger) {
	// Create a circuit breaker for HTTP client
	config := cbreak.DefaultConfig("http-client-example")
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

	// Create an HTTP client
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	// Simulate HTTP requests to failing endpoints
	endpoints := []string{
		"http://localhost:8080/error",
		"http://localhost:8080/timeout",
		"http://localhost:8080/unavailable",
	}

	log.Info("Simulating HTTP requests to failing endpoints...")
	for _, endpoint := range endpoints {
		result, err := breaker.Execute(context.Background(), func() (string, error) {
			resp, err := client.Get(endpoint)
			if err != nil {
				return "", err
			}
			defer resp.Body.Close()

			if resp.StatusCode >= 500 {
				return "", fmt.Errorf("server error: %d", resp.StatusCode)
			}
			return fmt.Sprintf("Status: %d", resp.StatusCode), nil
		})

		if err != nil {
			log.Error("Request to %s failed: %v", endpoint, err)
		} else {
			log.Success("Request to %s succeeded: %s", endpoint, result)
		}

		// Log the current state
		state := breaker.GetState()
		log.Info("Circuit breaker state: %s", state)
	}

	// Wait for the circuit breaker to reset
	log.Info("Waiting for circuit breaker to reset...")
	time.Sleep(6 * time.Second)

	// Try a request to a working endpoint
	log.Info("Trying a request to a working endpoint...")
	result, err := breaker.Execute(context.Background(), func() (string, error) {
		resp, err := client.Get("http://localhost:8080/health")
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 500 {
			return "", fmt.Errorf("server error: %d", resp.StatusCode)
		}
		return fmt.Sprintf("Status: %d", resp.StatusCode), nil
	})

	if err != nil {
		log.Error("Request failed: %v", err)
	} else {
		log.Success("Request succeeded: %s", result)
	}

	// Log the final state
	state := breaker.GetState()
	log.Info("Final circuit breaker state: %s", state)
}
