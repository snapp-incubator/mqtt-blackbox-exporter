package config

type Config struct {
	Trace  Trace  `json:"trace"  koanf:"trace"`
	Metric Metric `json:"metric" koanf:"metric"`
}

type Trace struct {
	Enabled  bool    `json:"enabled,omitempty"  koanf:"enabled"`
	Ratio    float64 `json:"ratio,omitempty"    koanf:"ratio"`
	Endpoint string  `json:"endpoint,omitempty" koanf:"endpoint"`
}

type Metric struct {
	Address string `json:"address,omitempty" koanf:"address"`
	Enabled bool   `json:"enabled,omitempty" koanf:"enabled"`
}
