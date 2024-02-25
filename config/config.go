package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type (
	Config struct {
		HTTP   `yaml:"http"`
		Logger `yaml:"logger"`
	}

	HTTP struct {
		Port             string `yaml:"port"`
		CorsAllowOrigins string `yaml:"cors_allow_origins"`
	}

	Logger struct {
		LogLevel string `yaml:"log_level"`
		Format   string `yaml:"format"`
		Color    bool   `yaml:"color"`
	}
)

func MustLoad() *Config {
	var cfg Config

	configPath := fetchConfigPath()

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}
	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
