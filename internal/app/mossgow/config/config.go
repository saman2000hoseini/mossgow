package config

import (
	"github.com/saman2000hoseini/mossgow/pkg/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

const (
	app       = "mossgow"
	cfgFile   = "config.yaml"
	cfgPrefix = "mossgow"
)

type (
	Config struct {
		PathLayers int      `mapstructure:"path-layers"`
		OutputDir  string   `mapstructure:"output-dir"`
		InputDir   string   `mapstructure:"input-dir"`
		MossDir    string   `mapstructure:"moss-dir"`
		Supported  []string `mapstructure:"supported"`
	}
)

func (c Config) Validate() error {
	return validator.New().Struct(c)
}

// Init initializes application configuration.
func Init() Config {
	var cfg Config

	config.Init(app, cfgFile, &cfg, defaultConfig, cfgPrefix)

	if err := cfg.Validate(); err != nil {
		logrus.Fatalf("failed to validate configurations: %s", err.Error())
	}

	return cfg
}
