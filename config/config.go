package config

import (
	"github.com/gresio/print-server/internal/auth"
	"github.com/gresio/print-server/internal/printer"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Printers map[string]printer.PrinterConfig `mapstructure:"printers"`
		Presets  map[string]Preset                `mapstructure:"presets"`
		HTTP     HttpConfig                       `mapstructure:"http"`
		Log      LogConfig                        `mapstructure:"logger"`
	}

	Preset struct {
		Printer       string
		JobAttributes printer.JobAttributes `mapstructure:",squash"`
	}

	HttpConfig struct {
		Port  string               `mapstructure:"port"`
		Users map[string]auth.User `mapstructure:"users"`
	}

	LogConfig struct {
		Level string
	}
)

func NewConfig() (*Config, error) {
	var cfg Config

	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AddConfigPath("../../config")
	v.AddConfigPath("./config")
	v.AddConfigPath("etc/cloudprint")
	v.ReadInConfig()

	err := v.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
