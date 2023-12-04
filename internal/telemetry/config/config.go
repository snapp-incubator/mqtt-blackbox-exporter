package config

type Config struct {
	Trace  Trace  `koanf:"trace"`
	Metric Metric `koanf:"metric"`
}

type Trace struct {
	Enabled  bool    `koanf:"enabled"`
	Ratio    float64 `koanf:"ratio"`
	Endpoint string  `koanf:"endpoint"`
}

type Metric struct {
	Address string `koanf:"address"`
	Enabled bool   `koanf:"enabled"`
}
