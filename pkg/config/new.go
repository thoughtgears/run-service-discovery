package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port         string `envconfig:"PORT" default:"8080"`
	GcpProjectID string `envconfig:"GCP_PROJECT_ID" required:"true"`
}

func NewConfig() (Config, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return Config{}, fmt.Errorf("failed to process env var: %w", err)
	}

	return c, nil
}
