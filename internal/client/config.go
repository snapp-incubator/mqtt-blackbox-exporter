package client

import (
	"time"
)

type Config struct {
	URL           string        `json:"url,omitempty"            koanf:"url"`
	ClientID      string        `json:"client_id,omitempty"      koanf:"clientid"`
	Username      string        `json:"username,omitempty"       koanf:"username"`
	Password      string        `json:"password,omitempty"       koanf:"password"`
	KeepAlive     time.Duration `json:"keep_alive,omitempty"     koanf:"keepalive"`
	PingTimeout   time.Duration `json:"ping_timeout,omitempty"   koanf:"ping_timeout"`
	AutoReconnect bool          `json:"auto_reconnect,omitempty" koanf:"auto_reconnect"`
	QoS           int           `json:"qo_s,omitempty"           koanf:"qos"`
	Retained      bool          `json:"retained,omitempty"       koanf:"retained"`
}
