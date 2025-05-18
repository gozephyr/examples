package main

import (
	"time"

	"github.com/gozephyr/examples/pkg/logger"
	"github.com/gozephyr/gencache"
	"github.com/gozephyr/gencache/policy"
)

func main() {
	log := logger.Get()
	log.SetPrefix("gencache-policy ")

	log.Section("LRU Policy Example")
	lruExample(log)

	log.Section("LFU Policy Example")
	lfuExample(log)

	log.Section("FIFO Policy Example")
	fifoExample(log)
}

func lruExample(log *logger.Logger) {
	// Create a cache with LRU policy
	cache := gencache.New[string, string](
		gencache.WithPolicy[string, string](policy.NewLRU[string, string](policy.WithMaxSize(3))),
	)
	defer func() {
		if err := cache.Close(); err != nil {
			log.Error("Error closing cache: %v", err)
		}
	}()

	// Add items to demonstrate LRU eviction
	items := []struct {
		key   string
		value string
	}{
		{"key1", "value1"},
		{"key2", "value2"},
		{"key3", "value3"},
		{"key4", "value4"}, // This should evict key1
	}

	log.Info("Adding items to cache...")
	for _, item := range items {
		err := cache.Set(item.key, item.value, time.Minute)
		if err != nil {
			log.Error("Error setting %s: %v", item.key, err)
			continue
		}
		log.Success("Added %s = %s", item.key, item.value)
	}

	// Access key2 to make it most recently used
	log.Info("Accessing key2 to make it most recently used...")
	value, err := cache.Get("key2")
	if err != nil {
		log.Error("Error getting key2: %v", err)
	} else {
		log.Success("Retrieved key2 = %s", value)
	}

	// Add another item, which should evict key3 (least recently used)
	log.Info("Adding key5, which should evict key3...")
	err = cache.Set("key5", "value5", time.Minute)
	if err != nil {
		log.Error("Error setting key5: %v", err)
	} else {
		log.Success("Added key5 = value5")
	}

	// Verify the state of the cache
	log.Info("Verifying cache state...")
	keys := []string{"key1", "key2", "key3", "key4", "key5"}
	for _, key := range keys {
		value, err := cache.Get(key)
		if err != nil {
			log.Success("Key %s was evicted as expected", key)
		} else {
			log.Info("Key %s is present with value: %s", key, value)
		}
	}
}

func lfuExample(log *logger.Logger) {
	// Create a cache with LFU policy
	cache := gencache.New[string, string](
		gencache.WithPolicy[string, string](policy.NewLFU[string, string](policy.WithMaxSize(3))),
	)
	defer func() {
		if err := cache.Close(); err != nil {
			log.Error("Error closing cache: %v", err)
		}
	}()

	// Add initial items
	items := []struct {
		key   string
		value string
	}{
		{"key1", "value1"},
		{"key2", "value2"},
		{"key3", "value3"},
	}

	log.Info("Adding initial items to cache...")
	for _, item := range items {
		err := cache.Set(item.key, item.value, time.Minute)
		if err != nil {
			log.Error("Error setting %s: %v", item.key, err)
			continue
		}
		log.Success("Added %s = %s", item.key, item.value)
	}

	// Access key1 multiple times to increase its frequency
	log.Info("Accessing key1 multiple times to increase its frequency...")
	for i := 0; i < 3; i++ {
		value, err := cache.Get("key1")
		if err != nil {
			log.Error("Error getting key1: %v", err)
		} else {
			log.Success("Retrieved key1 = %s (access %d)", value, i+1)
		}
	}

	// Access key2 once
	log.Info("Accessing key2 once...")
	value, err := cache.Get("key2")
	if err != nil {
		log.Error("Error getting key2: %v", err)
	} else {
		log.Success("Retrieved key2 = %s", value)
	}

	// Add a new item, which should evict key3 (least frequently used)
	log.Info("Adding key4, which should evict key3 (least frequently used)...")
	err = cache.Set("key4", "value4", time.Minute)
	if err != nil {
		log.Error("Error setting key4: %v", err)
	} else {
		log.Success("Added key4 = value4")
	}

	// Verify the state of the cache
	log.Info("Verifying cache state...")
	keys := []string{"key1", "key2", "key3", "key4"}
	for _, key := range keys {
		value, err := cache.Get(key)
		if err != nil {
			log.Success("Key %s was evicted as expected", key)
		} else {
			log.Info("Key %s is present with value: %s", key, value)
		}
	}
}

func fifoExample(log *logger.Logger) {
	// Create a cache with FIFO policy
	cache := gencache.New[string, string](
		gencache.WithPolicy[string, string](policy.NewFIFO[string, string](policy.WithMaxSize(3))),
	)
	defer func() {
		if err := cache.Close(); err != nil {
			log.Error("Error closing cache: %v", err)
		}
	}()

	// Add items to demonstrate FIFO eviction
	items := []struct {
		key   string
		value string
	}{
		{"key1", "value1"},
		{"key2", "value2"},
		{"key3", "value3"},
		{"key4", "value4"}, // This should evict key1
	}

	log.Info("Adding items to cache...")
	for _, item := range items {
		err := cache.Set(item.key, item.value, time.Minute)
		if err != nil {
			log.Error("Error setting %s: %v", item.key, err)
			continue
		}
		log.Success("Added %s = %s", item.key, item.value)
	}

	// Access key2 (should not affect eviction order)
	log.Info("Accessing key2 (should not affect eviction order)...")
	value, err := cache.Get("key2")
	if err != nil {
		log.Error("Error getting key2: %v", err)
	} else {
		log.Success("Retrieved key2 = %s", value)
	}

	// Add another item, which should evict key2 (oldest remaining item)
	log.Info("Adding key5, which should evict key2...")
	err = cache.Set("key5", "value5", time.Minute)
	if err != nil {
		log.Error("Error setting key5: %v", err)
	} else {
		log.Success("Added key5 = value5")
	}

	// Verify the state of the cache
	log.Info("Verifying cache state...")
	keys := []string{"key1", "key2", "key3", "key4", "key5"}
	for _, key := range keys {
		value, err := cache.Get(key)
		if err != nil {
			log.Success("Key %s was evicted as expected", key)
		} else {
			log.Info("Key %s is present with value: %s", key, value)
		}
	}
}
