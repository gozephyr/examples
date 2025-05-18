package main

import (
	"time"

	"github.com/gozephyr/examples/pkg/logger"
	"github.com/gozephyr/gencache"
)

// User represents a sample custom type for pooling
type User struct {
	ID   int
	Name string
}

func main() {
	log := logger.Get()
	log.SetPrefix("gencache-pooling ")
	log.Section("Object Pooling Example")
	poolingExample(log)
}

func poolingExample(log *logger.Logger) {
	// Create a cache with object pooling
	cache := gencache.New[string, *User](
		gencache.WithPoolConfig[string, *User](gencache.PoolConfig{
			MaxSize:       1000,
			MinSize:       10,
			CleanupPeriod: 5 * time.Minute,
			MaxIdleTime:   10 * time.Minute,
		}),
	)
	defer func() {
		if err := cache.Close(); err != nil {
			log.Error("Error closing cache: %v", err)
		}
	}()

	// Create and store users
	users := []*User{
		{ID: 1, Name: "John Doe"},
		{ID: 2, Name: "Jane Smith"},
		{ID: 3, Name: "Bob Johnson"},
	}

	// Store users in the pool
	for _, user := range users {
		key := user.Name
		log.Info("Storing user in pool: %+v", user)
		err := cache.Set(key, user, time.Minute)
		if err != nil {
			log.Error("Error storing user: %v", err)
			continue
		}
	}

	// Retrieve users from the pool
	for _, user := range users {
		key := user.Name
		log.Info("Retrieving user from pool: %s", key)
		retrievedUser, err := cache.Get(key)
		if err != nil {
			log.Error("Error retrieving user: %v", err)
			continue
		}
		log.Success("Retrieved user: ID=%d, Name=%s", retrievedUser.ID, retrievedUser.Name)
	}

	// Demonstrate pool reuse
	log.Info("Demonstrating pool reuse...")
	for i := 0; i < 3; i++ {
		user := &User{ID: i + 4, Name: "Pool User"}
		key := user.Name
		err := cache.Set(key, user, time.Minute)
		if err != nil {
			log.Error("Error storing user: %v", err)
			continue
		}
		log.Success("Stored user in pool: %+v", user)
	}

	// Demonstrate pool cleanup
	log.Info("Demonstrating pool cleanup...")
	if err := cache.Clear(); err != nil {
		log.Error("Error clearing pool: %v", err)
	} else {
		log.Success("Pool cleared successfully")
	}
}
