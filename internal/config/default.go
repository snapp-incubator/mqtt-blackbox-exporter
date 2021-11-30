package config

import (
	"github.com/snapp-incubator/mqtt-blackbox-exporter/internal/logger"
	telemetry "github.com/snapp-incubator/mqtt-blackbox-exporter/internal/telemetry/config"
)

// Default return default configuration.
func Default() Config {
	return Config{
		Logger: logger.Config{
			Level: "debug",
		},
		Telemetry: telemetry.Config{
			Trace: telemetry.Trace{
				Enabled: false,
				Ratio:   1.0,
				Agent: telemetry.Agent{
					Host: "127.0.0.1",
					Port: "6831",
				},
			},
			Profiler: telemetry.Profiler{
				Enabled: false,
				Address: "http://127.0.0.1:4040",
			},
			Metric: telemetry.Metric{
				Address: ":8080",
				Enabled: true,
			},
		},
	}
}