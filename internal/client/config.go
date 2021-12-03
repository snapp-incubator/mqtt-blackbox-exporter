package client

import (
	"time"
)

type Config struct {
	URL           string        `koanf:"url"`
	ClientID      string        `koanf:"clientid"`
	Username      string        `koanf:"username"`
	Password      string        `koanf:"password"`
	KeepAlive     time.Duration `koanf:"keepalive"`
	PingTimeout   time.Duration `koanf:"ping_timeout"`
	AutoReconnect bool          `koanf:"auto_reconnect"`
	QoS           int           `koanf:"qos"`
	Retained      bool          `koanf:"retained"`
}
