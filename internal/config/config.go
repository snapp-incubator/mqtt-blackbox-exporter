package config

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/snapp-incubator/mqtt-blackbox-exporter/internal/client"
	"github.com/snapp-incubator/mqtt-blackbox-exporter/internal/logger"
	telemetry "github.com/snapp-incubator/mqtt-blackbox-exporter/internal/telemetry/config"
	"github.com/tidwall/pretty"
)

const (
	// Prefix indicates environment variables prefix.
	Prefix = "mbe_"
)

type (
	// Config holds all configurations.
	Config struct {
		Telemetry    telemetry.Config `json:"telemetry,omitempty"     koanf:"telemetry"`
		Logger       logger.Config    `json:"logger,omitempty"        koanf:"logger"`
		MQTT         client.Config    `json:"mqtt,omitempty"          koanf:"mqtt"`
		PingDuration time.Duration    `json:"ping_duration,omitempty" koanf:"ping_duration"`
	}
)

// New reads configuration with koanf.
func New() Config {
	var instance Config

	k := koanf.New(".")

	// load default configuration from file
	err := k.Load(structs.Provider(Default(), "koanf"), nil)
	if err != nil {
		log.Fatalf("error loading default: %s", err)
	}

	// load configuration from file
	err = k.Load(file.Provider("config.yml"), yaml.Parser())
	if err != nil {
		log.Printf("error loading config.yml: %s", err)
	}

	// load environment variables
	err = k.Load(env.Provider(Prefix, ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, Prefix)), "__", ".")
	}), nil)
	if err != nil {
		log.Printf("error loading environment variables: %s", err)
	}

	err = k.Unmarshal("", &instance)
	if err != nil {
		log.Fatalf("error unmarshalling config: %s", err)
	}

	indent, err := json.MarshalIndent(instance, "", "\t")
	if err != nil {
		log.Fatalf("error marshaling config to json: %s", err)
	}

	indent = pretty.Color(indent, nil)
	tmpl := `
	================ Loaded Configuration ================
	%s
	======================================================
	`
	log.Printf(tmpl, string(indent))

	return instance
}
