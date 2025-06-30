package cmd

import (
	"os"

	"github.com/snapp-incubator/mqtt-blackbox-exporter/internal/config"
	"github.com/snapp-incubator/mqtt-blackbox-exporter/internal/logger"
	"github.com/snapp-incubator/mqtt-blackbox-exporter/internal/telemetry/metric"
	"github.com/snapp-incubator/mqtt-blackbox-exporter/internal/telemetry/trace"
	"github.com/spf13/cobra"
)

// ExitFailure status code.
const ExitFailure = 1

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cfg := config.New()

	logger := logger.New(cfg.Logger)

	tracer := trace.New(cfg.Telemetry.Trace)
	metric.NewServer(cfg.Telemetry.Metric).Start(logger.Named("metric"))

	// nolint: exhaustruct
	root := &cobra.Command{
		Use:   "mqtt-blackbox-exporter",
		Short: "ping pong with mqtt broker",
		Run: func(_ *cobra.Command, _ []string) {
			main(cfg, logger, tracer)
		},
	}

	err := root.Execute()
	if err != nil {
		os.Exit(ExitFailure)
	}
}
