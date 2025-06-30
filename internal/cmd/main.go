package cmd

import (
	"context"
	"time"

	"github.com/snapp-incubator/mqtt-blackbox-exporter/internal/client"
	"github.com/snapp-incubator/mqtt-blackbox-exporter/internal/config"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func main(cfg config.Config, logger *zap.Logger, trace trace.Tracer) {
	ctx := context.Background()

	{
		_, span := trace.Start(ctx, "cmd.main.subscriber")

		client := client.New(ctx, cfg.MQTT, logger.Named("mqtt"), trace, true)

		err := client.Connect(ctx)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			span.End()

			logger.Fatal("mqtt connection failed", zap.Error(err))
		}

		span.End()
	}

	{
		_, span := trace.Start(ctx, "cmd.main.publisher")

		client := client.New(ctx, cfg.MQTT, logger.Named("mqtt"), trace, false)

		err := client.Connect(ctx)
		if err != nil {
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

			err := client.Ping(ctx, int(i.UnixMilli()))
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())

				logger.Error("publish failed", zap.Error(err))
			}

			span.End()
		}
	}
}
