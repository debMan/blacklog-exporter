package config

import (
	"github.com/debman/blacklog-exporter/internal/kafkaclient"
	"github.com/debman/blacklog-exporter/internal/logger"
	"github.com/debman/blacklog-exporter/internal/metric"
)

// Default return default configuration.
// nolint: mnd
func Default() Config {
	return Config{
		Logger: logger.Config{
			Level: "info",
		},
		Kafka: kafkaclient.Config{
			BootstrapServers: []string{"localhost:9092"},
			AutoOffsetReset:  "earliest",
			Topics:           []string{"blacklogs"},
			GroupId:          "blacklog-exporter",
			AutoCommit:       true,
		},
		Metric: metric.Config{
			Server: metric.Server{
				Address: ":8080",
				Path:    "/metrics",
			},
			Enabled: true,
		},
	}
}
