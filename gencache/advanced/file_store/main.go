package main

import (
	"context"
	"time"

	"github.com/gozephyr/examples/pkg/logger"
	"github.com/gozephyr/gencache"
	"github.com/gozephyr/gencache/store"
)

func main() {
	log := logger.Get()
	log.SetPrefix("gencache-file ")
	log.Section("File Store Example")
	fileStoreExample(log)
}

func fileStoreExample(log *logger.Logger) {
	ctx := context.Background()

	// Create a file store
	config := &store.FileConfig{
		Directory:          "/tmp/gencache",
		FileExtension:      ".cache",
		CompressionEnabled: false,
		CompressionLevel:   0,
		CleanupInterval:    time.Hour,
	}
	fileStore, err := store.NewFileStore[string, string](ctx, config)
	if err != nil {
		log.Error("Error creating file store: %v", err)
		return
	}
	defer func() {
		if err := fileStore.Close(ctx); err != nil {
			log.Error("Error closing file store: %v", err)
		}
	}()

	// Create a cache with file store
	cache := gencache.New[string, string](
		gencache.WithStore[string, string](fileStore),
	)
	defer func() {
		if err := cache.Close(); err != nil {
			log.Error("Error closing cache: %v", err)
		}
	}()

	// Store some values
	values := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	log.Info("Storing values in file store...")
	for key, value := range values {
		err := cache.Set(key, value, time.Minute)
		if err != nil {
			log.Error("Error storing %s: %v", key, err)
			continue
		}
		log.Success("Stored %s = %s", key, value)
	}

	// Retrieve values
	log.Info("Retrieving values from file store...")
	for key := range values {
		value, err := cache.Get(key)
		if err != nil {
			log.Error("Error retrieving %s: %v", key, err)
			continue
		}
		log.Success("Retrieved %s = %s", key, value)
	}

	// Demonstrate persistence
	log.Info("Demonstrating persistence...")
	log.Info("Closing and reopening cache...")
	if err := cache.Close(); err != nil {
		log.Error("Error closing cache: %v", err)
	}
	if err := fileStore.Close(ctx); err != nil {
		log.Error("Error closing file store: %v", err)
	}

	// Reopen the store and cache
	fileStore, err = store.NewFileStore[string, string](ctx, config)
	if err != nil {
		log.Error("Error reopening file store: %v", err)
		return
	}
	defer func() {
		if err := fileStore.Close(ctx); err != nil {
			log.Error("Error closing file store: %v", err)
		}
	}()

	cache = gencache.New[string, string](
		gencache.WithStore[string, string](fileStore),
	)
	defer func() {
		if err := cache.Close(); err != nil {
			log.Error("Error closing cache: %v", err)
		}
	}()

	// Verify values are still there
	log.Info("Verifying persistence...")
	for key := range values {
		value, err := cache.Get(key)
		if err != nil {
			log.Error("Error retrieving %s after reopen: %v", key, err)
			continue
		}
		log.Success("Retrieved %s = %s after reopen", key, value)
	}

	// Demonstrate clear operation
	log.Info("Demonstrating clear operation...")
	if err := fileStore.Clear(ctx); err != nil {
		log.Error("Error clearing store: %v", err)
	} else {
		log.Success("Store cleared successfully")
	}
}
