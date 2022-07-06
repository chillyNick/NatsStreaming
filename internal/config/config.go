package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var cfg *Config

// GetConfigInstance returns service config
func GetConfigInstance() Config {
	if cfg != nil {
		return *cfg
	}

	return Config{}
}

type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type NatsStreaming struct {
	Port        int    `yaml:"port"`
	Host        string `yaml:"host"`
	ClusterId   string `yaml:"clusterId"`
	PublisherId string `yaml:"publisherId"`
	ConsumerId  string `yaml:"consumerId"`
	Subject     string `yaml:"subject"`
	DurableName string `yaml:"durableName"`
}

// Config - contains all configuration parameters in config package.
type Config struct {
	Database      Database      `yaml:"database"`
	NatsStreaming NatsStreaming `yaml:"natsStreaming"`
}

// ReadConfigYML - read configurations from file and init instance Config.
func ReadConfigYML(filePath string) error {
	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return err
	}

	return nil
}
