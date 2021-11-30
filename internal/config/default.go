package config

import telemetry "github.com/snapp-incubator/mqtt-blackbox-exporter/telemetry/config"

// Default return default configuration.
func Default() Config {
	return Config{
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
