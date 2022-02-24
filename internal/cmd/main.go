package cmd

import (
	"github.com/snapp-incubator/mqtt-blackbox-exporter/internal/cache"
	"time"

	"github.com/snapp-incubator/mqtt-blackbox-exporter/internal/client"
	"github.com/snapp-incubator/mqtt-blackbox-exporter/internal/config"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func main(cfg config.Config, logger *zap.Logger, trace trace.Tracer) {
	c := cache.Cache{}
	c.Init()

	{
		client := client.New(cfg.MQTT, logger.Named("mqtt"), trace, &c, true)

		if err := client.Connect(); err != nil {
			logger.Fatal("mqtt connection failed", zap.Error(err))
		}
	}

	{
		client := client.New(cfg.MQTT, logger.Named("mqtt"), trace, &c, false)

		if err := client.Connect(); err != nil {
			logger.Fatal("mqtt connection failed", zap.Error(err))
		}

		ticker := time.NewTicker(cfg.PingDuration)
		for i := range ticker.C {
			if err := client.Ping(c.Push(int(i.UnixMilli()), time.Now())); err != nil {
				logger.Error("publish failed", zap.Error(err))
			}
		}
	}
}
