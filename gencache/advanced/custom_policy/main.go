package main

import (
	"time"

	"github.com/gozephyr/examples/pkg/logger"
	"github.com/gozephyr/gencache"
)

// CustomPolicy implements a simple FIFO (First In, First Out) eviction policy
type CustomPolicy struct {
	items    []string
	capacity int
}

func NewCustomPolicy(capacity int) *CustomPolicy {
	return &CustomPolicy{
		items:    make([]string, 0, capacity),
		capacity: capacity,
	}
}

func (p *CustomPolicy) OnGet(key string, value string) {
	// No-op for FIFO policy
}

func (p *CustomPolicy) OnSet(key string, value string, ttl time.Duration) {
	if len(p.items) >= p.capacity {
		// Remove oldest item if at capacity
		p.items = p.items[1:]
	}
	p.items = append(p.items, key)
}

func (p *CustomPolicy) OnDelete(key string) {
	// Remove key from items
	for i, k := range p.items {
		if k == key {
			p.items = append(p.items[:i], p.items[i+1:]...)
			break
		}
	}
}

func (p *CustomPolicy) OnClear() {
	p.items = make([]string, 0, p.capacity)
}

func (p *CustomPolicy) Evict() (string, bool) {
	if len(p.items) == 0 {
		return "", false
	}
	key := p.items[0]
	p.items = p.items[1:]
	return key, true
}

func (p *CustomPolicy) Size() int {
	return len(p.items)
}

func (p *CustomPolicy) Capacity() int {
	return p.capacity
}

func main() {
	log := logger.Get()
	log.SetPrefix("gencache-policy ")
	log.Section("Custom Policy Example")
	customPolicyExample(log)
}

func customPolicyExample(log *logger.Logger) {
	// Create a custom policy with capacity of 3
	customPolicy := NewCustomPolicy(3)

	// Create a cache with the custom policy
	cache := gencache.New[string, string](
		gencache.WithPolicy[string, string](customPolicy),
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

	// Verify eviction
	log.Info("Verifying eviction...")
	for _, item := range items {
		value, err := cache.Get(item.key)
		if err != nil {
			log.Success("Key %s was evicted as expected", item.key)
		} else {
			log.Info("Key %s is present with value: %s", item.key, value)
		}
	}

	// Demonstrate policy behavior
	log.Section("Policy Statistics")
	log.Info("Current Size: %d", customPolicy.Size())
	log.Info("Capacity: %d", customPolicy.Capacity())
	log.Info("Items in order: %v", customPolicy.items)
}
