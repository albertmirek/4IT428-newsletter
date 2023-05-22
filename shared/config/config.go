package config

import (
	"github.com/kelseyhightower/envconfig"
	postgresql "vse.com/4IT428/2023/newsletter/shared/db/posgtresql"
)

type Config struct {
	Address  string `envconfig:"ADDRESS" default:":8080"`
	Database postgresql.Config
}

func LoadConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
