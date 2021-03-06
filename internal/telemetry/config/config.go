package config

type Config struct {
	Trace    Trace    `koanf:"trace"`
	Profiler Profiler `koanf:"profiler"`
	Metric   Metric   `koanf:"metric"`
}

type Trace struct {
	Enabled bool `koanf:"enabled"`
	Agent   `koanf:"agent"`
	Ratio   float64 `koanf:"ratio"`
}

type Agent struct {
	Host string `koanf:"host"`
	Port string `koanf:"port"`
}

type Profiler struct {
	Enabled bool   `koanf:"enabled"`
	Address string `koanf:"address"`
}

type Metric struct {
	Address string `koanf:"address"`
	Enabled bool   `koanf:"enabled"`
}
