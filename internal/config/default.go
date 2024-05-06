package config

import (
	"time"

	"github.com/snapp-incubator/mqtt-blackbox-exporter/internal/client"
	"github.com/snapp-incubator/mqtt-blackbox-exporter/internal/logger"
	telemetry "github.com/snapp-incubator/mqtt-blackbox-exporter/internal/telemetry/config"
)

// Default return default configuration.
// nolint: mnd
func Default() Config {
	return Config{
		Telemetry: telemetry.Config{
			Trace: telemetry.Trace{
				Enabled:  false,
				Ratio:    0.1,
				Endpoint: "127.0.0.1:4317",
			},
			Metric: telemetry.Metric{
				Address: ":8080",
				Enabled: true,
			},
		},
		Logger: logger.Config{
			Level: "debug",
		},
		MQTT: client.Config{
			URL:           "tcp://127.0.0.1:1883",
			ClientID:      "",
			Username:      "",
			Password:      "",
			KeepAlive:     60 * time.Second,
			PingTimeout:   1 * time.Second,
			AutoReconnect: true,
			QoS:           1,
			Retained:      true,
		},
		PingDuration: 1 * time.Second,
	}
}
