package client

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	PingTopic = "snapp/ping"

	DisconnectTimeout = 250
)

// Client wraps mqtt client for handling publishing and subscribing.
type Client struct {
	Client  mqtt.Client
	Logger  *zap.Logger
	Tracer  trace.Tracer
	Metrics Metrics

	QoS         int
	Retained    bool
	IsSubscribe bool
}

// Message contains the information to send over ping.
type Message struct {
	Headers map[string]string `json:"headers,omitempty"`
}

// New creates a new mqtt client with given configuration.
// isSubscribe for ping message.
// nolint: funlen
func New(ctx context.Context,
	cfg Config,
	logger *zap.Logger,
	tracer trace.Tracer,
	isSubscribe bool,
) *Client {
	mqtt.DEBUG, _ = zap.NewStdLogAt(logger.Named("raw"), zap.DebugLevel)
	mqtt.ERROR, _ = zap.NewStdLogAt(logger.Named("raw"), zap.ErrorLevel)

	_, span := tracer.Start(ctx, "client.new")
	defer span.End()

	clientID := cfg.ClientID
	if clientID == "" {
		var err error

		clientID, err = os.Hostname()
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			logger.Fatal("hostname fetching failed, specify a client id", zap.Error(err))
		}
	}

	if isSubscribe {
		clientID += "-subscriber"
	} else {
		clientID += "-producer"
	}

	span.SetAttributes(attribute.String("client-id", clientID))
	span.SetAttributes(attribute.String("broker-url", cfg.URL))

	opts := mqtt.NewClientOptions()

	opts.SetClientID(clientID)
	opts.AddBroker(cfg.URL)

	if cfg.Username != "" {
		opts.SetUsername(cfg.Username)
	}

	if cfg.Password != "" {
		opts.SetPassword(cfg.Password)
	}

	opts.SetKeepAlive(cfg.KeepAlive)
	opts.SetAutoReconnect(cfg.AutoReconnect)
	opts.SetPingTimeout(cfg.PingTimeout)

	client := &Client{
		Logger:      logger,
		Client:      nil,
		Tracer:      tracer,
		Metrics:     NewMetrics(),
		Retained:    cfg.Retained,
		QoS:         cfg.QoS,
		IsSubscribe: isSubscribe,
	}

	opts.SetConnectionLostHandler(client.OnConnectionLostHandler)
	opts.SetOnConnectHandler(client.OnConnectHandler)

	client.Client = mqtt.NewClient(opts)

	return client
}

func (c *Client) OnConnectionLostHandler(_ mqtt.Client, err error) {
	c.Logger.Error("connection lost", zap.Error(err))
	c.Metrics.ConnectionErrors.Add(1)
}

func (c *Client) OnConnectHandler(_ mqtt.Client) {
	ctx := context.Background()
	_, span := c.Tracer.Start(ctx, "client.on.connect.handler")

	defer span.End()

	if c.IsSubscribe {
		if token := c.Client.Subscribe(PingTopic, byte(c.QoS), c.Pong); token.Wait() && token.Error() != nil {
			span.RecordError(token.Error())
			span.SetStatus(codes.Error, token.Error().Error())

			c.Logger.Fatal("subscription failed", zap.String("topic", PingTopic), zap.Error(token.Error()))
		}
	}
}

func (c *Client) Pong(_ mqtt.Client, b mqtt.Message) {
	var msg Message

	ctx := otel.GetTextMapPropagator().Extract(context.Background(), propagation.MapCarrier(msg.Headers))

	_, span := c.Tracer.Start(ctx, "ping.received", trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()

	err := json.Unmarshal(b.Payload(), &msg)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		c.Logger.Fatal("cannot marshal message", zap.Error(err))
	}

	if id, has := msg.Headers["id"]; has {
		if value, has := msg.Headers["start"]; has {
			start, _ := strconv.Atoi(value)
			duration := time.Since(time.UnixMilli(int64(start)))

			c.Logger.Info("successful ping", zap.Duration("time", duration), zap.String("id", id))
			c.Metrics.PingDuration.Observe(duration.Seconds())
			c.Metrics.Pongs.Inc()
		}
	}
}

func (c *Client) Ping(ctx context.Context, id int) error {
	c.Logger.Debug("ping...", zap.String("topic", PingTopic))

	var msg Message

	msg.Headers = make(map[string]string)
	msg.Headers["id"] = strconv.Itoa(id)
	msg.Headers["start"] = strconv.Itoa(int(time.Now().UnixMilli()))

	_, span := c.Tracer.Start(ctx, "ping.publish", trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(msg.Headers))

	b, err := json.Marshal(msg)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		c.Logger.Fatal("cannot marshal message", zap.Error(err))
	}

	if token := c.Client.Publish(PingTopic, byte(c.QoS), c.Retained, b); token.Wait() && token.Error() != nil {
		span.RecordError(token.Error())
		span.SetStatus(codes.Error, token.Error().Error())

		c.Metrics.PublishErrors.Inc()

		return fmt.Errorf("failed to publish %w", err)
	}

	c.Metrics.Pings.Inc()

	return nil
}

func (c *Client) Disconnect() {
	c.Client.Disconnect(DisconnectTimeout)
}

func (c *Client) Connect(ctx context.Context) error {
	_, span := c.Tracer.Start(ctx, "client.on.connect")

	defer span.End()

	if token := c.Client.Connect(); token.Wait() && token.Error() != nil {
		c.Metrics.ConnectionErrors.Add(1)

		span.SetStatus(codes.Error, token.Error().Error())
		span.RecordError(token.Error())

		return fmt.Errorf("mqtt connection failed %w", token.Error())
	}

	return nil
}
