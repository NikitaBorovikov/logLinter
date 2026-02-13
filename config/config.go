package config

import "github.com/ilyakaznacheev/cleanenv"

const (
	configYAMLPath = "config/config.yaml"
)

type Config struct {
	SensitiveWords []string `yaml:"sensitive-words"`
	LogPkgNames    []string `yaml:"log-package-names"`
	LogIndentNames []string `yaml:"log-indent-names"`
}

func InitConfig() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(configYAMLPath, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
