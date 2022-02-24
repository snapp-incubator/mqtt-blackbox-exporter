package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/snapp-incubator/mqtt-blackbox-exporter/internal/cache"
	"os"
	"strconv"
	"time"

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
	Cache   *cache.Cache

	QoS         int
	Retained    bool
	IsSubscribe bool
}

// Message contains the information to send over ping.
type Message struct {
	Headers map[string]string
}

// New creates a new mqtt client with given configuration.
// isSubscribe for ping message.
func New(cfg Config, logger *zap.Logger, tracer trace.Tracer, cache *cache.Cache, isSubscribe bool) *Client {
	mqtt.DEBUG, _ = zap.NewStdLogAt(logger.Named("raw"), zap.DebugLevel)
	mqtt.ERROR, _ = zap.NewStdLogAt(logger.Named("raw"), zap.ErrorLevel)

	clientID := cfg.ClientID
	if clientID == "" {
		var err error

		clientID, err = os.Hostname()
		if err != nil {
			logger.Fatal("hostname fetching failed, specify a client id", zap.Error(err))
		}
	}

	if isSubscribe {
		clientID += "-subscriber"
	} else {
		clientID += "-producer"
	}

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
		Cache:       cache,
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
	if c.IsSubscribe {
		if token := c.Client.Subscribe(PingTopic, byte(c.QoS), c.Pong); token.Wait() && token.Error() != nil {
			c.Logger.Fatal("subscription failed", zap.String("topic", PingTopic), zap.Error(token.Error()))
		}
	}
}

func (c *Client) Pong(_ mqtt.Client, b mqtt.Message) {
	var msg Message

	if err := json.Unmarshal(b.Payload(), &msg); err != nil {
		c.Logger.Fatal("cannot marshal message", zap.Error(err))
	}

	ctx := otel.GetTextMapPropagator().Extract(context.Background(), propagation.MapCarrier(msg.Headers))

	_, span := c.Tracer.Start(ctx, "ping.received")
	defer span.End()

	if value, has := msg.Headers["id"]; has {
		id, _ := strconv.Atoi(value)
		item := c.Cache.Pull(id)

		if item.Status {
			item.Status = false
			duration := time.Since(item.Start)

			c.Logger.Info("successful ping", zap.Duration("time", duration), zap.Int("id", id))
			c.Metrics.PingDuration.Observe(duration.Seconds())
		}
	}
}

func (c *Client) Ping(id int) error {
	c.Logger.Debug("ping...", zap.String("topic", PingTopic))

	var msg Message
	msg.Headers = make(map[string]string)
	msg.Headers["id"] = strconv.Itoa(id)

	ctx, span := c.Tracer.Start(context.Background(), "ping.publish")
	defer span.End()

	otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(msg.Headers))

	b, err := json.Marshal(msg)
	if err != nil {
		c.Logger.Fatal("cannot marshal message", zap.Error(err))
	}

	if token := c.Client.Publish(PingTopic, byte(c.QoS), c.Retained, b); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to publish %w", err)
	}

	return nil
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
