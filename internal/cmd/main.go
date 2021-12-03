package cmd

import (
	"github.com/snapp-incubator/mqtt-blackbox-exporter/internal/client"
	"github.com/snapp-incubator/mqtt-blackbox-exporter/internal/config"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func main(cfg config.Config, logger *zap.Logger, trace trace.Tracer) {
	client := client.New(cfg.MQTT, logger.Named("mqtt"))

	if err := client.Connect(); err != nil {
		logger.Fatal("mqtt connection failed", zap.Error(err))
	}
}
