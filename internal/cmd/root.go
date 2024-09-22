package cmd

import (
	"os"

	"github.com/debman/blacklog-exporter/internal/config"
	"github.com/debman/blacklog-exporter/internal/logger"
	"github.com/debman/blacklog-exporter/internal/metric"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

// ExitFailure status code.
const ExitFailure = 1

func Execute() {
	var configPath string
	flag.StringVarP(&configPath, "config", "c", "./config.yaml", "Path to config file")
	flag.Parse()

	cfg := config.New(configPath)

	logger := logger.New(cfg.Logger)

	metric.NewServer(cfg.Metric).Start(logger.Named("metric"))

	// nolint: exhaustruct
	root := &cobra.Command{
		Use:   "blacklog-exporter",
		Short: "ping pong with kafka broker",
		Run: func(_ *cobra.Command, _ []string) {
			main(cfg, logger)
		},
	}

	if err := root.Execute(); err != nil {
		os.Exit(ExitFailure)
	}
}
