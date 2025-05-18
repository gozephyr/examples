package main

import (
	"context"
	"time"

	"github.com/gozephyr/examples/pkg/logger"
	"github.com/gozephyr/gencache"
)

func main() {
	log := logger.Get()
	log.SetPrefix("gencache-batch ")
	log.Section("Batch Operations Example")
	batchExample(log)
}

func batchExample(log *logger.Logger) {
	ctx := context.Background()

	// Create a cache with batch operations
	cache := gencache.New[string, string](
		gencache.WithBatchConfig[string, string](gencache.BatchConfig{
			MaxBatchSize:     1000,
			OperationTimeout: 5 * time.Second,
			MaxConcurrent:    10,
		}),
	)
	defer func() {
		if err := cache.Close(); err != nil {
			log.Error("Error closing cache: %v", err)
		}
	}()

	// Create batch cache
	batchCache := gencache.NewBatchCache(cache, gencache.DefaultBatchConfig())

	// Prepare batch data
	keys := []string{"key1", "key2", "key3", "key4", "key5"}
	values := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
		"key4": "value4",
		"key5": "value5",
	}

	// Batch set operation
	log.Info("Performing batch set operation...")
	err := batchCache.SetMany(ctx, values, time.Minute)
	if err != nil {
		log.Error("Batch set error: %v", err)
		return
	}
	log.Success("Batch set completed successfully")

	// Batch get operation
	log.Info("Performing batch get operation...")
	retrievedValues := batchCache.GetMany(ctx, keys)
	for key, value := range retrievedValues {
		if value != "" {
			log.Success("Retrieved %s = %s", key, value)
		} else {
			log.Warn("No value found for key: %s", key)
		}
	}

	// Batch delete operation
	log.Info("Performing batch delete operation...")
	err = batchCache.DeleteMany(ctx, keys[:3]) // Delete first 3 keys
	if err != nil {
		log.Error("Batch delete error: %v", err)
		return
	}
	log.Success("Batch delete completed successfully")

	// Verify deletion
	log.Info("Verifying deleted keys...")
	for _, key := range keys[:3] {
		value, err := cache.Get(key)
		if err != nil {
			log.Success("Key %s was successfully deleted", key)
		} else {
			log.Warn("Key %s still exists with value: %s", key, value)
		}
	}
}
