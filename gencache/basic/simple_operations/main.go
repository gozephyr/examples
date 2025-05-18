package main

import (
	"time"

	"github.com/gozephyr/examples/pkg/logger"
	"github.com/gozephyr/gencache"
)

// User represents a sample custom type for caching
type User struct {
	ID   int
	Name string
}

func main() {
	log := logger.Get()
	log.SetPrefix("gencache ")

	log.Section("Basic Cache Operations Example")
	basicOperations(log)

	log.Section("Custom Type Operations Example")
	customTypeOperations(log)

	log.Section("TTL Operations Example")
	ttlOperations(log)
}

func basicOperations(log *logger.Logger) {
	cache := gencache.New[string, string]()

	defer func() {
		if err := cache.Close(); err != nil {
			log.Warn("Error closing cache: %v", err)
		}
	}()

	log.Info("Setting value 'value1' for key 'key1'...")
	err := cache.Set("key1", "value1", time.Minute)
	if err != nil {
		log.Error("Error setting value: %v", err)
		return
	}
	log.Success("Value set successfully")

	log.Info("Getting value for key 'key1'...")
	value, err := cache.Get("key1")
	if err != nil {
		log.Error("Error getting value: %v", err)
		return
	}
	log.Success("Retrieved value: %s", value)

	log.Info("Deleting key 'key1'...")
	err = cache.Delete("key1")
	if err != nil {
		log.Error("Error deleting value: %v", err)
		return
	}
	log.Success("Key deleted successfully")

	log.Info("Trying to get deleted key 'key1'...")
	_, err = cache.Get("key1")
	if err != nil {
		log.Warn("Expected error for deleted key: %v", err)
	}
}

func customTypeOperations(log *logger.Logger) {
	cache := gencache.New[string, *User]()

	defer func() {
		if err := cache.Close(); err != nil {
			log.Warn("Error closing cache: %v", err)
		}
	}()

	user := &User{ID: 1, Name: "John Doe"}
	log.Info("Storing user: %+v", user)
	err := cache.Set("user1", user, time.Minute)
	if err != nil {
		log.Error("Error setting user: %v", err)
		return
	}
	log.Success("User stored successfully")

	log.Info("Retrieving user with key 'user1'...")
	retrievedUser, err := cache.Get("user1")
	if err != nil {
		log.Error("Error getting user: %v", err)
		return
	}
	log.Success("Retrieved user: ID=%d, Name=%s", retrievedUser.ID, retrievedUser.Name)
}

func ttlOperations(log *logger.Logger) {
	cache := gencache.New[string, string]()

	defer func() {
		if err := cache.Close(); err != nil {
			log.Warn("Error closing cache: %v", err)
		}
	}()

	log.Info("Setting value with 1s TTL...")
	err := cache.Set("temp_key", "temp_value", time.Second)
	if err != nil {
		log.Error("Error setting value: %v", err)
		return
	}
	log.Success("Value set successfully")

	log.Info("Getting value immediately...")
	value, err := cache.Get("temp_key")
	if err != nil {
		log.Error("Error getting value: %v", err)
		return
	}
	log.Success("Retrieved value: %s", value)

	log.Info("Waiting for value to expire (2s)...")
	time.Sleep(2 * time.Second)

	log.Info("Trying to get expired value...")
	_, err = cache.Get("temp_key")
	if err != nil {
		log.Warn("Expected error for expired key: %v", err)
	}
}
