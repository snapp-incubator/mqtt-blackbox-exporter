package cmd

import (
	"time"

	"github.com/snapp-incubator/mqtt-blackbox-exporter/internal/client"
	"github.com/snapp-incubator/mqtt-blackbox-exporter/internal/config"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func main(cfg config.Config, logger *zap.Logger, trace trace.Tracer) {
	{
		client := client.New(cfg.MQTT, logger.Named("mqtt"), trace, true)

		if err := client.Connect(); err != nil {
			logger.Fatal("mqtt connection failed", zap.Error(err))
		}
	}

	{
		client := client.New(cfg.MQTT, logger.Named("mqtt"), trace, false)

		if err := client.Connect(); err != nil {
			logger.Fatal("mqtt connection failed", zap.Error(err))
		}

		ticker := time.NewTicker(cfg.PingDuration)
		for range ticker.C {
			client.Ping()
		}
	}
}
