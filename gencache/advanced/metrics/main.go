package main

import (
	"time"

	"github.com/gozephyr/examples/pkg/logger"
	"github.com/gozephyr/gencache"
	"github.com/gozephyr/gencache/metrics"
)

func main() {
	log := logger.Get()
	log.SetPrefix("gencache-metrics ")
	log.Section("Metrics Example")
	metricsExample(log)
	prometheusMetricsExample(log)
}

func metricsExample(log *logger.Logger) {
	// Create a cache with metrics using the built-in exporter
	cache := gencache.New[string, string](
		gencache.WithMetricsConfig[string, string](gencache.MetricsConfig{
			ExporterType: metrics.StandardExporter,
			CacheName:    "gencache-examples",
			Labels: map[string]string{
				"environment": "production",
			},
		}),
	)
	defer func() {
		if err := cache.Close(); err != nil {
			log.Error("Error closing cache: %v", err)
		}
	}()

	// Perform some operations to generate metrics
	operations := []struct {
		op    string
		key   string
		value string
	}{
		{"set", "key1", "value1"},
		{"set", "key2", "value2"},
		{"get", "key1", ""},
		{"get", "key2", ""},
		{"get", "key3", ""}, // This will be a miss
		{"delete", "key1", ""},
		{"get", "key1", ""}, // This will be a miss after deletion
	}

	for _, op := range operations {
		var err error
		switch op.op {
		case "set":
			err = cache.Set(op.key, op.value, time.Minute)
			if err != nil {
				log.Error("Set error: %v", err)
			} else {
				log.Success("Set %s = %s", op.key, op.value)
			}
		case "get":
			value, err := cache.Get(op.key)
			if err != nil {
				log.Warn("Get miss for %s: %v", op.key, err)
			} else {
				log.Success("Get hit for %s = %s", op.key, value)
			}
		case "delete":
			err = cache.Delete(op.key)
			if err != nil {
				log.Error("Delete error: %v", err)
			} else {
				log.Success("Deleted %s", op.key)
			}
		}
	}

	// Metrics are exported via the configured exporter (e.g., metrics.StandardExporter)
	// Check your application logs or the appropriate endpoint for metrics output.
	log.Section("Metrics Output")
	log.Info("Metrics are exported via the configured exporter. Check logs or the appropriate endpoint for output.")
}

// prometheusMetricsExample demonstrates how to use Prometheus metrics with gencache
func prometheusMetricsExample(log *logger.Logger) {
	log.Section("Prometheus Metrics Example")
	// Create a cache with Prometheus metrics enabled
	cache := gencache.New[string, string](
		gencache.WithMetricsConfig[string, string](gencache.MetricsConfig{
			ExporterType: metrics.PrometheusExporterType,
			CacheName:    "gencache-prometheus-example",
			Labels: map[string]string{
				"environment": "production",
			},
		}),
	)
	defer func() {
		if err := cache.Close(); err != nil {
			log.Error("Error closing cache: %v", err)
		}
	}()

	// Perform some operations to generate metrics
	operations := []struct {
		op    string
		key   string
		value string
	}{
		{"set", "key1", "value1"},
		{"set", "key2", "value2"},
		{"get", "key1", ""},
		{"get", "key2", ""},
		{"get", "key3", ""}, // This will be a miss
		{"delete", "key1", ""},
		{"get", "key1", ""}, // This will be a miss after deletion
	}

	for _, op := range operations {
		var err error
		switch op.op {
		case "set":
			err = cache.Set(op.key, op.value, time.Minute)
			if err != nil {
				log.Error("Set error: %v", err)
			} else {
				log.Success("Set %s = %s", op.key, op.value)
			}
		case "get":
			value, err := cache.Get(op.key)
			if err != nil {
				log.Warn("Get miss for %s: %v", op.key, err)
			} else {
				log.Success("Get hit for %s = %s", op.key, value)
			}
		case "delete":
			err = cache.Delete(op.key)
			if err != nil {
				log.Error("Delete error: %v", err)
			} else {
				log.Success("Deleted %s", op.key)
			}
		}
	}

	// Metrics are exported via Prometheus
	// Check the /metrics endpoint for Prometheus metrics output
	log.Section("Prometheus Metrics Output")
	log.Info("Metrics are exported via Prometheus. Check the /metrics endpoint for output.")
}
