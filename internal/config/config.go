package config

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/debman/blacklog-exporter/internal/kafkaclient"
	"github.com/debman/blacklog-exporter/internal/logger"
	"github.com/debman/blacklog-exporter/internal/metric"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/tidwall/pretty"
)

const (
	// Prefix indicates environment variables prefix.
	Prefix = "BLE_"
)

type (
	// Config holds all configurations.
	Config struct {
		Logger logger.Config      `json:"logger,omitempty" koanf:"logger"`
		Kafka  kafkaclient.Config `json:"kafka,omitempty"   koanf:"kafka"`
		Metric metric.Config      `json:"metric,omitempty" koanf:"metric"`
	}
)

// New reads configuration with koanf.
func New(configPath string) Config {
	var instance Config

	k := koanf.New(".")

	// load default configuration from file
	if err := k.Load(structs.Provider(Default(), "koanf"), nil); err != nil {
		log.Fatalf("error loading default: %s", err)
	}

	// load configuration from file
	if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
		log.Printf("error loading config.yaml from: %s", configPath)
	}

	// load environment variables
	if err := k.Load(env.Provider(Prefix, ".", func(s string) string {
		placeholder := "\x01"
		tmpStr := strings.ReplaceAll(s, "__", placeholder)
		dottedStr := strings.ReplaceAll(tmpStr, "_", ".")
		result := strings.ReplaceAll(dottedStr, placeholder, "_")
		return strings.ToLower(result)
	}), nil); err != nil {
		log.Printf("error loading environment variables: %s", err)
	}

	if err := k.Unmarshal("", &instance); err != nil {
		log.Fatalf("error unmarshalling config: %s", err)
	}

	if instance.Logger.Level == "debug" {
		indent, err := json.MarshalIndent(instance, "", "  ")
		if err != nil {
			log.Fatalf("error marshalling config to json: %s", err)
		}
		indent = pretty.Color(indent, nil)
		configOut := fmt.Sprintf("\n==== Configuration:\n%s\n====", indent)
		log.Print(configOut)
	}

	return instance
}
