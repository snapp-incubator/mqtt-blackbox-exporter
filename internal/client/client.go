package client

import (
	"context"
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.opentelemetry.io/otel"
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

	QoS      int
	Retained bool
}

// Message contains the information to send over ping.
type Message struct {
	Headers map[string]string
}

// New creates a new mqtt client with given configuration.
// isSubscribe for ping message.
func New(cfg Config, logger *zap.Logger, tracer trace.Tracer, isSubscribe bool) *Client {
	mqtt.DEBUG, _ = zap.NewStdLogAt(logger.Named("raw"), zap.DebugLevel)
	mqtt.ERROR, _ = zap.NewStdLogAt(logger.Named("raw"), zap.ErrorLevel)

	opts := mqtt.NewClientOptions()

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
		Logger:   logger,
		Client:   mqtt.NewClient(opts),
		Tracer:   tracer,
		Metrics:  NewMetrics(),
		Retained: cfg.Retained,
		QoS:      cfg.QoS,
	}

	opts.SetConnectionLostHandler(client.OnConnectionLostHandler)

	return client
}

func (c *Client) OnConnectionLostHandler(_ mqtt.Client, err error) {
	c.Logger.Error("connection lost", zap.Error(err))
	c.Metrics.ConnectionErrors.Add(1)
}

func (c *Client) OnConnectHandler(_ mqtt.Client) {
	if token := c.Client.Subscribe(PingTopic, byte(c.QoS), c.Pong); token.Wait() && token.Error() != nil {
		c.Logger.Fatal("subscription failed", zap.String("topic", PingTopic), zap.Error(token.Error()))
	}
}

func (c *Client) Pong(_ mqtt.Client, b mqtt.Message) {
	ctx, span := c.Tracer.Start(context.Background(), "ping.subscribe")
	defer span.End()

	var msg Message

	if err := json.Unmarshal(b.Payload(), &msg); err != nil {
		c.Logger.Fatal("cannot marshal message", zap.Error(err))
	}

	otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(msg.Headers))
}

func (c *Client) Ping() {
	var msg Message

	ctx, span := c.Tracer.Start(context.Background(), "ping.publish")
	defer span.End()

	otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(msg.Headers))

	b, err := json.Marshal(msg)
	if err != nil {
		c.Logger.Fatal("cannot marshal message", zap.Error(err))
	}

	c.Client.Publish(PingTopic, byte(c.QoS), c.Retained, b)
}

func (c *Client) Disconnect() {
	c.Client.Disconnect(DisconnectTimeout)
}

func (c *Client) Connect() error {
	if token := c.Client.Connect(); token.Wait() && token.Error() != nil {
		c.Metrics.ConnectionErrors.Add(1)

		return fmt.Errorf("mqtt connection failed %w", token.Error())
	}

	return nil
}
