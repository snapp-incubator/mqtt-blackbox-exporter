package client

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
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
	Metrics Metrics
}

// New creates a new mqtt client with given configuration.
func New(cfg Config, logger *zap.Logger) *Client {
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

	client := mqtt.NewClient(opts)

	return &Client{
		Logger:  logger,
		Client:  client,
		Metrics: NewMetrics(),
	}
}

func (c *Client) OnConnectionLostHandler(_ mqtt.Client, err error) {
	c.Logger.Error("connection lost", zap.Error(err))
	c.Metrics.ConnectionErrors.Add(1)
}

func (c *Client) OnConnectHandler() {
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
