package main

import (
	"time"

	"github.com/gozephyr/examples/pkg/logger"
	"github.com/gozephyr/gencache"
)

func main() {
	log := logger.Get()
	log.SetPrefix("gencache-error ")

	log.Section("Error Handling Example")
	errorHandlingExample(log)
}

func errorHandlingExample(log *logger.Logger) {
	cache := gencache.New[string, string]()
	defer func() {
		if err := cache.Close(); err != nil {
			log.Warn("Error closing cache: %v", err)
		}
	}()

	// Try to get a non-existent key
	log.Info("Trying to get a non-existent key 'missing_key'...")
	_, err := cache.Get("missing_key")
	if err != nil {
		log.Error("Get error: %v", err)
	}

	// Try to delete a non-existent key
	log.Info("Trying to delete a non-existent key 'missing_key'...")
	err = cache.Delete("missing_key")
	if err != nil {
		log.Error("Delete error: %v", err)
	}

	// Try to set a key with too short TTL
	log.Info("Setting key 'short_lived' with 1s TTL...")
	err = cache.Set("short_lived", "value", time.Second)
	if err != nil {
		log.Error("Set error: %v", err)
	}

	// Wait for the key to expire
	log.Info("Waiting for 2s to let the key expire...")
	time.Sleep(2 * time.Second)

	// Try to get the expired key
	_, err = cache.Get("short_lived")
	if err != nil {
		log.Error("Get after expiration: %v", err)
	}
}
