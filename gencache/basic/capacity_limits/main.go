package main

import (
	"fmt"
	"time"

	"github.com/gozephyr/examples/pkg/logger"
	"github.com/gozephyr/gencache"
)

func main() {
	log := logger.Get()
	log.SetPrefix("gencache-capacity")
	log.Section("Cache Operations Example")
	cacheOperationsExample(log)
}

func cacheOperationsExample(log *logger.Logger) {
	// Create a cache
	cache := gencache.New[string, string]()
	defer func() {
		if err := cache.Close(); err != nil {
			log.Warn("Error closing cache: %v", err)
		}
	}()

	// Insert items
	for i := 1; i <= 4; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		log.Info("Setting %s = %s", key, value)
		err := cache.Set(key, value, time.Minute)
		if err != nil {
			log.Error("Set error: %v", err)
		}
	}

	// Check which keys are present
	for i := 1; i <= 4; i++ {
		key := fmt.Sprintf("key%d", i)
		value, err := cache.Get(key)
		if err != nil {
			log.Warn("%s was not found: %v", key, err)
		} else {
			log.Success("%s is present with value: %s", key, value)
		}
	}
}
