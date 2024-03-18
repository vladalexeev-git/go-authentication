package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type (
	Config struct {
		HTTP        `yaml:"http"`
		Logger      `yaml:"logger"`
		Postgres    `yaml:"postgres"`
		AccessToken `yaml:"access_token"`
		Session     `yaml:"session"`
		MongoDB     `yaml:"mongodb"`
	}

	HTTP struct {
		Port             string `yaml:"port"`
		CorsAllowOrigins string `yaml:"cors_allow_origins"`
	}
	Logger struct {
		Env string `yaml:"env"`
	}
	Postgres struct {
		PoolMax int    `yaml:"pool_max"`
		URL     string `env:"PG_URL"`
	}
	Session struct {
		TTL            time.Duration `yaml:"ttl"`
		CookieKey      string        `yaml:"cookie_key"`
		CookieDomain   string        `yaml:"cookie_domain"`
		CookiePath     string        `yaml:"cookie_path"`
		CookieSecure   bool          `yaml:"cookie_secure"`
		CookieHttpOnly bool          `yaml:"cookie_httponly"`
	}
	AccessToken struct {
		TTL        time.Duration `yaml:"ttl"`
		SigningKey string        `yaml:"signing_key"`
	}

	MongoDB struct {
		DbName   string `yaml:"db_name"`
		URI      string `yaml:"uri" env:"MONGO_URI"`
		Username string `yaml:"username" env:"MONGO_USER"`
		Password string `yaml:"password" env:"MONGO_PASS"`
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
