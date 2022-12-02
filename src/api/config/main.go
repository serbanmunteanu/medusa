package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type WorkerServerConfig struct {
	RedisConfig    RedisConfig `yaml:"redisConfig"`
	UseCustomQueue bool        `yaml:"useCustomQueue"`
	Queues         []Queue     `yaml:"queues"`
	Concurrency    int         `yaml:"concurrency"`
	StrictPriority bool        `yaml:"strictPriority"`
}

type WebServerConfig struct {
}

type Queue struct {
	Name     string `yaml:"name"`
	Priority int    `yaml:"priority"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
}

func LoadConfig(configType string, config interface{}) {
	var configPath string

	switch configType {
	case "worker":
		configPath = "api/config/worker-config.yaml"
	case "web":
		configPath = "api/config/web-config.yaml"
	default:
		panic("config type not configured")
	}

	yamlConfig, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	parsingErr := yaml.Unmarshal(yamlConfig, config)
	if parsingErr != nil {
		panic(err)
	}
}
