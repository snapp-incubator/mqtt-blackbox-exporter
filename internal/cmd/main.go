package cmd

import (
	"context"
	"time"

	"github.com/snapp-incubator/mqtt-blackbox-exporter/internal/cache"
	"github.com/snapp-incubator/mqtt-blackbox-exporter/internal/client"
	"github.com/snapp-incubator/mqtt-blackbox-exporter/internal/config"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func main(cfg config.Config, logger *zap.Logger, trace trace.Tracer) {
	c := cache.Cache{}
	c.Init()

	{
		_, span := trace.Start(context.Background(), "cmd.main.subscriber")
		defer span.End()

		client := client.New(cfg.MQTT, logger.Named("mqtt"), trace, &c, true)

		if err := client.Connect(); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			logger.Fatal("mqtt connection failed", zap.Error(err))
		}
	}

	{
		_, span := trace.Start(context.Background(), "cmd.main.publisher")
		defer span.End()

		client := client.New(cfg.MQTT, logger.Named("mqtt"), trace, &c, false)

		if err := client.Connect(); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			logger.Fatal("mqtt connection failed", zap.Error(err))
		}

		ticker := time.NewTicker(cfg.PingDuration)
		for i := range ticker.C {
			if err := client.Ping(int(i.UnixMilli())); err != nil {
				span.SetAttributes(attribute.Int("ping_id", int(i.UnixMilli())))
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())

				logger.Error("publish failed", zap.Error(err))
			}
		}
	}
}
