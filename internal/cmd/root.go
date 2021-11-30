package cmd

import (
	"os"

	"github.com/snapp-incubator/mqtt-blackbox-exporter/internal/config"
	"github.com/snapp-incubator/mqtt-blackbox-exporter/internal/telemetry/trace"
	"github.com/spf13/cobra"
)

// ExitFailure status code.
const ExitFailure = 1

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cfg := config.New()

	_ = trace.New(cfg.Telemetry.Trace)

	// nolint: exhaustivestruct
	root := &cobra.Command{
		Use:   "mqtt-blackbox-exporter",
		Short: "ping pong with mqtt broker",
	}

	if err := root.Execute(); err != nil {
		os.Exit(ExitFailure)
	}
}
