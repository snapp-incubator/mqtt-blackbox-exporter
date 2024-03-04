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
	ctx := context.Background()

	c := cache.Cache{}
	c.Init()

	{
		_, span := trace.Start(ctx, "cmd.main.subscriber")

		client := client.New(ctx, cfg.MQTT, logger.Named("mqtt"), trace, &c, true)

		if err := client.Connect(ctx); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			span.End()

			logger.Fatal("mqtt connection failed", zap.Error(err))
		}

		span.End()
	}

	{
		_, span := trace.Start(ctx, "cmd.main.publisher")

		client := client.New(ctx, cfg.MQTT, logger.Named("mqtt"), trace, &c, false)

		if err := client.Connect(ctx); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			span.End()

			logger.Fatal("mqtt connection failed", zap.Error(err))
		}

		span.End()

		ticker := time.NewTicker(cfg.PingDuration)
		for i := range ticker.C {
			_, span := trace.Start(ctx, "main.ping")
			span.SetAttributes(attribute.Int("ping_id", int(i.UnixMilli())))

			if err := client.Ping(ctx, int(i.UnixMilli())); err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())

				logger.Error("publish failed", zap.Error(err))
			}

			span.End()
		}
	}
}
