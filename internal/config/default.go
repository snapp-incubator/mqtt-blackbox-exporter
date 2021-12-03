package config

import (
	"time"

	"github.com/snapp-incubator/mqtt-blackbox-exporter/internal/client"
	"github.com/snapp-incubator/mqtt-blackbox-exporter/internal/logger"
	telemetry "github.com/snapp-incubator/mqtt-blackbox-exporter/internal/telemetry/config"
)

// Default return default configuration.
// nolint: gomnd
func Default() Config {
	return Config{
		Telemetry: telemetry.Config{
			Trace: telemetry.Trace{
				Enabled: false,
				Ratio:   0.1,
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
		PingDuration: 60 * time.Second,
	}
}
